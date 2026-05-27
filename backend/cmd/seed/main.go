package main

import (
	"fmt"
	"math/rand"

	"inference-engine/internal/article"
	"inference-engine/internal/comment"
	"inference-engine/internal/config"
	"inference-engine/internal/pkg/database"
	"inference-engine/internal/prompt"
	"inference-engine/internal/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	config.Init()
	cfg := config.Get()

	db := database.InitMySQL(&cfg.Database)

	fmt.Println("seeding data...")

	seed(db)

	fmt.Println("done!")
}

func seed(db *gorm.DB) {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	admin := user.User{Username: "admin", Email: "admin@test.com", Password: string(hashedPwd), Role: "admin", Bio: "站点管理员", Status: 1}
	alice := user.User{Username: "alice", Email: "alice@test.com", Password: string(hashedPwd), Role: "creator", Bio: "全栈工程师", Status: 1}
	bob := user.User{Username: "bob", Email: "bob@test.com", Password: string(hashedPwd), Role: "creator", Bio: "AI 爱好者", Status: 1}
	users := []user.User{admin, alice, bob}
	for _, u := range users {
		db.Where("email = ?", u.Email).FirstOrCreate(&u)
	}
	var u1, u2, u3 user.User
	db.Where("email = ?", "admin@test.com").First(&u1)
	db.Where("email = ?", "alice@test.com").First(&u2)
	db.Where("email = ?", "bob@test.com").First(&u3)

	categories := []article.Category{
		{Name: "Go", Desc: "Go 语言相关", Sort: 1},
		{Name: "AI", Desc: "人工智能与机器学习", Sort: 2},
		{Name: "前端", Desc: "前端开发技术", Sort: 3},
		{Name: "后端", Desc: "后端架构与设计", Sort: 4},
		{Name: "DevOps", Desc: "运维与部署", Sort: 5},
	}
	for _, c := range categories {
		db.Where("name = ?", c.Name).FirstOrCreate(&c)
	}
	var catGo, catAI, catFront, catBack, catDev article.Category
	db.Where("name = ?", "Go").First(&catGo)
	db.Where("name = ?", "AI").First(&catAI)
	db.Where("name = ?", "前端").First(&catFront)
	db.Where("name = ?", "后端").First(&catBack)
	db.Where("name = ?", "DevOps").First(&catDev)

	tags := []article.Tag{
		{Name: "gin", Color: "#6c5ce7"},
		{Name: "gorm", Color: "#00b894"},
		{Name: "docker", Color: "#0984e3"},
		{Name: "vue", Color: "#00cec9"},
		{Name: "gpt", Color: "#e17055"},
		{Name: "redis", Color: "#d63031"},
		{Name: "mysql", Color: "#fdcb6e"},
		{Name: "websocket", Color: "#6c5ce7"},
	}
	for _, t := range tags {
		db.Where("name = ?", t.Name).FirstOrCreate(&t)
	}

	articles := []article.Article{
		{UserID: u2.ID, Title: "Go 1.24 新特性详解", Content: "# Go 1.24 新特性\n\nGo 1.24 带来了许多令人兴奋的新特性...\n\n## 泛型类型别名\n\n```go\ntype MySlice[T any] = []T\n```\n\n## 性能改进\n\n内存分配和垃圾回收都有显著优化。\n\n## 总结\n\n升级到 Go 1.24 会带来更好的开发体验。", Summary: "Go 1.24 在泛型、性能和工具链方面都有重要更新，本文详细介绍每个新特性。", CategoryID: uintPtr(catGo.ID), Tags: `["go","gin"]`, Status: "published", ViewCount: 1280, LikeCount: 45, CommentCount: 3},
		{UserID: u3.ID, Title: "Prompt Engineering 最佳实践", Content: "# Prompt Engineering 指南\n\n写好 prompt 是使用 AI 的关键技能。\n\n## 明确目标\n\n每次对话前明确你想要的输出格式。\n\n## 提供上下文\n\n```markdown\n你是一个资深 Go 工程师，请帮我审查以下代码：\n```\n\n## 迭代优化\n\n不断调整 prompt 直到得到理想结果。", Summary: "掌握 prompt engineering 的核心技巧，让 AI 输出更准确。", CategoryID: uintPtr(catAI.ID), Tags: `["gpt"]`, Status: "published", ViewCount: 2560, LikeCount: 89, CommentCount: 5},
		{UserID: u2.ID, Title: "Vue 3 + TypeScript 项目搭建", Content: "# Vue 3 + TypeScript 项目搭建\n\n## 环境准备\n\n```bash\nnpm create vue@latest\n```\n\n## 目录结构\n\n```\nsrc/\n  components/\n  views/\n  stores/\n  router/\n```\n\n## 状态管理\n\n使用 Pinia 替代 Vuex，类型更安全。", Summary: "从零搭建 Vue 3 + TypeScript 企业级项目的完整指南。", CategoryID: uintPtr(catFront.ID), Tags: `["vue"]`, Status: "published", ViewCount: 890, LikeCount: 32, CommentCount: 2},
		{UserID: u3.ID, Title: "Docker Compose 多服务部署", Content: "# Docker Compose 实战\n\n## 配置示例\n\n```yaml\nversion: '3.8'\nservices:\n  app:\n    build: .\n    ports:\n      - \"8080:8080\"\n```\n\n## 网络配置\n\n使用自定义网络让容器间通信。\n\n## 数据持久化\n\n使用 named volume 保存数据库数据。", Summary: "使用 Docker Compose 编排多服务应用的最佳实践。", CategoryID: uintPtr(catDev.ID), Tags: `["docker"]`, Status: "published", ViewCount: 1560, LikeCount: 67, CommentCount: 4},
		{UserID: u1.ID, Title: "GORM 使用技巧汇总", Content: "# GORM 使用技巧\n\n## 软删除\n\n```go\ntype Model struct {\n  DeletedAt gorm.DeletedAt\n}\n```\n\n## 预加载\n\n```go\ndb.Preload(\"User\").Find(&articles)\n```\n\n## 事务\n\n```go\ndb.Transaction(func(tx *gorm.DB) error {\n  return nil\n})\n```", Summary: "日常开发中常用的 GORM 技巧和踩坑记录。", CategoryID: uintPtr(catGo.ID), Tags: `["go","gorm"]`, Status: "published", ViewCount: 720, LikeCount: 28, CommentCount: 1},
		{UserID: u2.ID, Title: "WebSocket 实时通信方案选型", Content: "# WebSocket 方案对比\n\n| 方案 | 优点 | 缺点 |\n|------|------|------|\n| Gorilla | 成熟稳定 | 需自行管理 |\n| nhooyr | 标准库风格 | 社区较小 |\n\n## 推荐\n\n中小项目推荐 Gorilla WebSocket。", Summary: "Go 语言 WebSocket 库的对比与选择建议。", CategoryID: uintPtr(catBack.ID), Tags: `["go","websocket"]`, Status: "published", ViewCount: 630, LikeCount: 22, CommentCount: 2},
	}
	for _, a := range articles {
		db.Where("title = ?", a.Title).FirstOrCreate(&a)
	}

	comments := []comment.Comment{
		{ArticleID: 1, UserID: u3.ID, Content: "写得很好，总结到位！"},
		{ArticleID: 1, UserID: u1.ID, Content: "补充一点，泛型性能在 1.24 有大幅提升。"},
		{ArticleID: 2, UserID: u2.ID, Content: "chain-of-thought 技巧也很实用。"},
		{ArticleID: 2, UserID: u1.ID, Content: "感谢分享，正好需要这个。"},
		{ArticleID: 3, UserID: u3.ID, Content: "Vite 比 Webpack 快太多了。"},
	}
	for _, c := range comments {
		db.Where("content = ?", c.Content).FirstOrCreate(&c)
	}

	prompts := []prompt.Prompt{
		{UserID: u2.ID, Title: "代码审查助手", Content: "你是一个资深的 {{language}} 工程师，请审查以下代码：\n\n```{{language}}\n{{code}}\n```\n\n请从以下方面给出建议：\n1. 代码质量\n2. 性能问题\n3. 安全漏洞\n4. 可维护性", Description: "让 AI 帮你审查代码，支持任意编程语言。", Category: "开发", Model: "gpt-4", Tags: `["代码审查","编程"]`, UsageCount: 89, LikeCount: 34, Rating: 4.5},
		{UserID: u3.ID, Title: "SQL 优化专家", Content: "请优化以下 SQL 查询：\n\n```sql\n{{sql}}\n```\n\n表结构：\n```sql\n{{schema}}\n```\n\n请提供：\n1. 问题分析\n2. 优化后的 SQL\n3. 索引建议", Description: "输入慢 SQL，获取优化建议和索引推荐。", Category: "数据库", Model: "gpt-4", Tags: `["SQL","数据库","优化"]`, UsageCount: 120, LikeCount: 56, Rating: 4.8},
		{UserID: u2.ID, Title: "API 文档生成", Content: "根据以下代码生成 OpenAPI 规范的 API 文档：\n\n```{{language}}\n{{code}}\n```\n\n输出格式为 YAML，包含：\n- 接口路径\n- 请求方法\n- 参数说明\n- 响应格式\n- 示例", Description: "自动生成 RESTful API 文档，支持 Swagger/OpenAPI 格式。", Category: "开发", Model: "gpt-3.5", Tags: `["API","文档","OpenAPI"]`, UsageCount: 45, LikeCount: 18, Rating: 4.2},
		{UserID: u3.ID, Title: "单元测试生成器", Content: "为以下 {{language}} 函数生成全面的单元测试：\n\n```{{language}}\n{{code}}\n```\n\n测试要求：\n- 覆盖正常路径\n- 覆盖边界条件\n- 覆盖错误处理\n- 使用 {{test_framework}} 框架", Description: "输入函数代码，自动生成单元测试用例。", Category: "测试", Model: "gpt-4", Tags: `["测试","单元测试","TDD"]`, UsageCount: 67, LikeCount: 41, Rating: 4.6},
		{UserID: u2.ID, Title: "架构设计评审", Content: "请评审以下系统架构：\n\n## 背景\n{{background}}\n\n## 技术栈\n{{stack}}\n\n## 架构方案\n{{architecture}}\n\n请从以下角度分析：\n1. 可扩展性\n2. 可用性\n3. 成本\n4. 运维复杂度", Description: "让 AI 帮你评审系统架构方案，发现潜在问题。", Category: "架构", Model: "gpt-4", Tags: `["架构","设计","评审"]`, UsageCount: 33, LikeCount: 15, Rating: 4.3},
	}
	for _, p := range prompts {
		db.Where("title = ?", p.Title).FirstOrCreate(&p)
	}

	// follows: alice关注bob, bob关注alice
	var follows []user.Follow
	uIDs := []uint{u2.ID, u3.ID}
	for _, id1 := range uIDs {
		for _, id2 := range uIDs {
			if id1 != id2 {
				follows = append(follows, user.Follow{FollowerID: id1, FollowedID: id2})
			}
		}
	}
	for _, f := range follows {
		db.Where("follower_id = ? AND followed_id = ?", f.FollowerID, f.FollowedID).FirstOrCreate(&f)
	}

	// update stat counts
	db.Model(&user.User{}).Where("id = ?", u2.ID).Updates(map[string]any{
		"follow_count": 1, "fans_count": 1, "article_count": 3,
	})
	db.Model(&user.User{}).Where("id = ?", u3.ID).Updates(map[string]any{
		"follow_count": 1, "fans_count": 1, "article_count": 2,
	})

	// random likes on articles
	var allArticles []article.Article
	db.Find(&allArticles)
	for _, a := range allArticles {
		for _, uid := range []uint{u2.ID, u3.ID} {
			if a.UserID != uid {
				db.Where("user_id = ? AND target_type = 'article' AND target_id = ?", uid, a.ID).FirstOrCreate(&article.Like{UserID: uid, TargetType: "article", TargetID: a.ID})
			}
		}
	}

	_ = rand.Intn(10)
	fmt.Println("seeded: users(3), categories(5), tags(8), articles(6), comments(5), prompts(5), follows(2), likes(10)")
}

func uintPtr(v uint) *uint {
	return &v
}
