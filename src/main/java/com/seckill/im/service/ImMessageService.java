package com.seckill.im.service;

import com.seckill.im.config.KafkaConfig;
import com.seckill.im.entity.*;
import com.seckill.im.mapper.*;
import com.seckill.im.proto.ImMessageProto;
import com.mybatisflex.core.query.QueryWrapper;
import com.seckill.im.server.ChannelAttr;
import com.seckill.im.server.ChannelManager;
import com.seckill.im.util.SnowflakeIdGenerator;
import io.netty.channel.Channel;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Set;
import java.util.concurrent.CompletableFuture;

@Slf4j
@Service
@RequiredArgsConstructor
public class ImMessageService {

    private final ImMessageMapper imMessageMapper;
    private final ImSessionMapper imSessionMapper;
    private final ImSessionUserMapper imSessionUserMapper;
    private final ImGroupMemberMapper imGroupMemberMapper;
    private final ImReadReceiptMapper imReadReceiptMapper;
    private final ImRedisService imRedisService;
    private final KafkaTemplate<String, byte[]> kafkaTemplate;
    private final SnowflakeIdGenerator snowflakeIdGenerator;

    public void processUpstreamMessage(Channel senderChannel, Long userId, ImMessageProto.UpstreamMessage msg) {
        if (!imRedisService.checkAndMarkIdempotent(msg.getClientMsgId())) {
            ImMessageProto.MessageAck ack = ImMessageProto.MessageAck.newBuilder()
                    .setClientMsgId(msg.getClientMsgId())
                    .setCode(2)
                    .setErrMsg("Duplicate message")
                    .build();
            sendToChannel(senderChannel, ImMessageProto.WrapperMessage.newBuilder()
                    .setType(ImMessageProto.WrapperMessage.MessageType.MESSAGE_ACK)
                    .setMessageAck(ack)
                    .build());
            return;
        }

        Long sessionId = msg.getSessionId();
        if (sessionId == null || sessionId == 0) {
            sessionId = resolveOrCreateSession(userId, msg.getToUserId(), msg.getGroupId());
        }

        long serverMsgId = snowflakeIdGenerator.nextId();
        long seq = imRedisService.generateSessionSeq(sessionId);

        ImMessage messageEntity = ImMessage.builder()
                .id(serverMsgId)
                .clientMsgId(msg.getClientMsgId())
                .sessionId(sessionId)
                .senderId(userId)
                .msgType(msg.getMsgType().getNumber())
                .content(msg.getContent())
                .mediaUrl(msg.getMediaUrl())
                .seq(seq)
                .status(1)
                .createTime(LocalDateTime.now())
                .build();
        imMessageMapper.insert(messageEntity);

        ImMessageProto.DownstreamMessage downstream = ImMessageProto.DownstreamMessage.newBuilder()
                .setClientMsgId(msg.getClientMsgId())
                .setServerMsgId(serverMsgId)
                .setSessionId(sessionId)
                .setSessionType(msg.getGroupId() != 0 ? ImMessageProto.SessionType.GROUP : ImMessageProto.SessionType.SINGLE)
                .setSenderId(userId)
                .setToUserId(msg.getToUserId())
                .setGroupId(msg.getGroupId())
                .setMsgType(msg.getMsgType())
                .setContent(msg.getContent())
                .setMediaUrl(msg.getMediaUrl())
                .setSeq(seq)
                .setStatus(ImMessageProto.MsgStatus.NORMAL)
                .setServerTime(System.currentTimeMillis())
                .build();

        ImMessageProto.MessageAck ack = ImMessageProto.MessageAck.newBuilder()
                .setClientMsgId(msg.getClientMsgId())
                .setServerMsgId(serverMsgId)
                .setServerTime(System.currentTimeMillis())
                .setCode(0)
                .build();
        sendToChannel(senderChannel, ImMessageProto.WrapperMessage.newBuilder()
                .setType(ImMessageProto.WrapperMessage.MessageType.MESSAGE_ACK)
                .setMessageAck(ack)
                .build());

        CompletableFuture<SendResult<String, byte[]>> future =
                kafkaTemplate.send(KafkaConfig.TOPIC_IM_MESSAGE,
                        String.valueOf(sessionId),
                        downstream.toByteArray());

        future.whenComplete((result, ex) -> {
            if (ex != null) {
                log.error("[ImMessageService] Kafka send failed, clientMsgId:{}", msg.getClientMsgId(), ex);
            } else {
                log.debug("[ImMessageService] Kafka sent, clientMsgId:{}, offset:{}",
                        msg.getClientMsgId(), result.getRecordMetadata().offset());
            }
        });
    }

    @KafkaListener(topics = KafkaConfig.TOPIC_IM_MESSAGE, groupId = "im-message-consumer", concurrency = "4")
    public void consumeImMessage(List<byte[]> messages) {
        for (byte[] data : messages) {
            try {
                ImMessageProto.DownstreamMessage msg = ImMessageProto.DownstreamMessage.parseFrom(data);
                dispatchToReceivers(msg);
            } catch (Exception e) {
                log.error("[ImMessageService] consume error", e);
            }
        }
    }

    private void dispatchToReceivers(ImMessageProto.DownstreamMessage msg) {
        ImMessageProto.WrapperMessage wrapper = ImMessageProto.WrapperMessage.newBuilder()
                .setType(ImMessageProto.WrapperMessage.MessageType.DOWNSTREAM_MESSAGE)
                .setDownstreamMessage(msg)
                .build();

        if (msg.getSessionType() == ImMessageProto.SessionType.SINGLE) {
            dispatchSingleChat(msg, wrapper);
        } else {
            dispatchGroupChat(msg, wrapper);
        }
    }

    private void dispatchSingleChat(ImMessageProto.DownstreamMessage msg, ImMessageProto.WrapperMessage wrapper) {
        Long toUserId = msg.getToUserId();
        Set<Channel> channels = ChannelManager.getUserChannels(toUserId);

        if (channels.isEmpty() || !imRedisService.isUserOnline(toUserId)) {
            log.info("[ImMessageService] user:{} offline, storing offline message", toUserId);
            imRedisService.addOfflineMessage(toUserId, msg);
            imRedisService.incrementUnreadCount(toUserId, msg.getSessionId());
            return;
        }

        for (Channel channel : channels) {
            sendToChannel(channel, wrapper);
        }

        imRedisService.incrementUnreadCount(toUserId, msg.getSessionId());
    }

    private void dispatchGroupChat(ImMessageProto.DownstreamMessage msg, ImMessageProto.WrapperMessage wrapper) {
        List<Long> memberIds = imGroupMemberMapper.selectMemberUserIdsByGroupId(msg.getGroupId());
        for (Long memberId : memberIds) {
            if (memberId.equals(msg.getSenderId())) {
                continue;
            }

            Set<Channel> channels = ChannelManager.getUserChannels(memberId);
            if (channels.isEmpty() || !imRedisService.isUserOnline(memberId)) {
                imRedisService.addOfflineMessage(memberId, msg);
                imRedisService.incrementUnreadCount(memberId, msg.getSessionId());
                continue;
            }

            for (Channel channel : channels) {
                sendToChannel(channel, wrapper);
            }
            imRedisService.incrementUnreadCount(memberId, msg.getSessionId());
        }
    }

    public void processReadReceipt(Long userId, ImMessageProto.ReadReceipt receipt) {
        imRedisService.updateLastReadSeq(userId, receipt.getSessionId(), receipt.getLastReadSeq());
        imRedisService.resetUnreadCount(userId, receipt.getSessionId());
        imSessionUserMapper.updateLastReadSeq(receipt.getSessionId(), userId, receipt.getLastReadSeq());

        kafkaTemplate.send(KafkaConfig.TOPIC_IM_READ_RECEIPT,
                String.valueOf(receipt.getSessionId()),
                receipt.toByteArray());
    }

    @Transactional
    public void processRecall(Long userId, ImMessageProto.RecallMessage recall) {
        int updated = imMessageMapper.recallMessage(recall.getMsgId(), userId);
        if (updated > 0) {
            ImMessageProto.WrapperMessage wrapper = ImMessageProto.WrapperMessage.newBuilder()
                    .setType(ImMessageProto.WrapperMessage.MessageType.RECALL_MESSAGE)
                    .setRecallMessage(recall)
                    .build();

            ImMessage msg = imMessageMapper.selectOneById(recall.getMsgId());
            if (msg != null) {
                kafkaTemplate.send(KafkaConfig.TOPIC_IM_RECALL,
                        String.valueOf(msg.getSessionId()),
                        recall.toByteArray());
            }
        }
    }

    public void processSyncOffline(Channel channel, Long userId, ImMessageProto.SyncOfflineRequest request) {
        List<ImMessageProto.DownstreamMessage> messages =
                imRedisService.getAndRemoveOfflineMessages(userId, 100);

        ImMessageProto.SyncOfflineResponse response = ImMessageProto.SyncOfflineResponse.newBuilder()
                .addAllMessages(messages)
                .setTotal(messages.size())
                .setHasMore(false)
                .build();

        sendToChannel(channel, ImMessageProto.WrapperMessage.newBuilder()
                .setType(ImMessageProto.WrapperMessage.MessageType.SYNC_OFFLINE_RESPONSE)
                .setSyncOfflineResponse(response)
                .build());
    }

    public void handleUserOffline(Long userId) {
        imRedisService.markUserOffline(userId);

        ImMessageProto.OnlineStatusChange statusChange = ImMessageProto.OnlineStatusChange.newBuilder()
                .setUserId(userId)
                .setOnlineStatus(0)
                .build();

        kafkaTemplate.send(KafkaConfig.TOPIC_IM_ONLINE_STATUS,
                String.valueOf(userId),
                statusChange.toByteArray());
    }

    @KafkaListener(topics = KafkaConfig.TOPIC_IM_RECALL, groupId = "im-recall-consumer", concurrency = "2")
    public void consumeRecall(List<byte[]> messages) {
        for (byte[] data : messages) {
            try {
                ImMessageProto.RecallMessage recall = ImMessageProto.RecallMessage.parseFrom(data);
                broadcastRecall(recall);
            } catch (Exception e) {
                log.error("[ImMessageService] consume recall error", e);
            }
        }
    }

    private void broadcastRecall(ImMessageProto.RecallMessage recall) {
        ImMessageProto.WrapperMessage wrapper = ImMessageProto.WrapperMessage.newBuilder()
                .setType(ImMessageProto.WrapperMessage.MessageType.RECALL_MESSAGE)
                .setRecallMessage(recall)
                .build();

        ImMessage msg = imMessageMapper.selectOneById(recall.getMsgId());
        if (msg == null) return;

        if (isGroupSession(msg.getSessionId())) {
            List<Long> memberIds = imGroupMemberMapper.selectMemberUserIdsByGroupId(
                    imSessionMapper.selectOneById(msg.getSessionId()).getGroupId());
            for (Long memberId : memberIds) {
                for (Channel ch : ChannelManager.getUserChannels(memberId)) {
                    sendToChannel(ch, wrapper);
                }
            }
        } else {
            ImSession session = imSessionMapper.selectOneById(msg.getSessionId());
            if (session != null) {
                List<ImSessionUser> sessionUsers = imSessionUserMapper
                        .selectListByQuery(new QueryWrapper()
                                .eq("session_id", msg.getSessionId()));
                for (ImSessionUser su : sessionUsers) {
                    for (Channel ch : ChannelManager.getUserChannels(su.getUserId())) {
                        sendToChannel(ch, wrapper);
                    }
                }
            }
        }
    }

    public int getUnreadCount(Long userId) {
        return imRedisService.getUnreadCount(userId);
    }

    private Long resolveOrCreateSession(Long userId, Long toUserId, Long groupId) {
        if (groupId != null && groupId != 0) {
            return groupId;
        }
        ImSession session = imSessionMapper.selectSingleSessionByUserIds(userId, toUserId);
        if (session != null) {
            return session.getId();
        }
        return createSingleSession(userId, toUserId);
    }

    @Transactional
    public Long createSingleSession(Long userId1, Long userId2) {
        ImSession session = ImSession.builder()
                .sessionType(1)
                .status(1)
                .createTime(LocalDateTime.now())
                .build();
        imSessionMapper.insert(session);

        ImSessionUser su1 = ImSessionUser.builder()
                .sessionId(session.getId())
                .userId(userId1)
                .joinedTime(LocalDateTime.now())
                .build();
        ImSessionUser su2 = ImSessionUser.builder()
                .sessionId(session.getId())
                .userId(userId2)
                .joinedTime(LocalDateTime.now())
                .build();
        imSessionUserMapper.insert(su1);
        imSessionUserMapper.insert(su2);

        return session.getId();
    }

    private boolean isGroupSession(Long sessionId) {
        ImSession session = imSessionMapper.selectOneById(sessionId);
        return session != null && session.getSessionType() == 2;
    }

    private void sendToChannel(Channel channel, ImMessageProto.WrapperMessage msg) {
        if (channel != null && channel.isActive() && channel.isWritable()) {
            channel.writeAndFlush(msg, channel.voidPromise());
        }
    }
}