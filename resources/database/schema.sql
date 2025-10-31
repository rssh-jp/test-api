CREATE DATABASE IF NOT EXISTS testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE testdb;

-- =====================================================
-- ユーザーテーブル（基本情報）
-- =====================================================
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT 'ユーザー名（ログインID）',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT 'メールアドレス',
    password_hash VARCHAR(255) NOT NULL COMMENT 'パスワードハッシュ',
    status ENUM('active', 'inactive', 'suspended', 'deleted') NOT NULL DEFAULT 'active' COMMENT 'アカウント状態',
    email_verified BOOLEAN DEFAULT FALSE COMMENT 'メール認証済みフラグ',
    last_login_at TIMESTAMP NULL COMMENT '最終ログイン日時',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_last_login (last_login_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザー基本情報';

-- =====================================================
-- ユーザープロフィールテーブル（詳細情報）
-- =====================================================
CREATE TABLE IF NOT EXISTS user_profiles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE COMMENT 'ユーザーID',
    first_name VARCHAR(50) COMMENT '名',
    last_name VARCHAR(50) COMMENT '姓',
    display_name VARCHAR(100) COMMENT '表示名',
    bio TEXT COMMENT '自己紹介',
    avatar_url VARCHAR(500) COMMENT 'アバター画像URL',
    birth_date DATE COMMENT '生年月日',
    gender ENUM('male', 'female', 'other', 'prefer_not_to_say') COMMENT '性別',
    country_code CHAR(2) COMMENT '国コード（ISO 3166-1 alpha-2）',
    timezone VARCHAR(50) DEFAULT 'UTC' COMMENT 'タイムゾーン',
    language_code CHAR(2) DEFAULT 'en' COMMENT '言語コード',
    phone_number VARCHAR(20) COMMENT '電話番号',
    website_url VARCHAR(500) COMMENT 'ウェブサイトURL',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_country (country_code),
    INDEX idx_display_name (display_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザープロフィール詳細';

-- =====================================================
-- カテゴリーテーブル
-- =====================================================
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE COMMENT 'カテゴリー名',
    slug VARCHAR(100) NOT NULL UNIQUE COMMENT 'URLスラッグ',
    description TEXT COMMENT '説明',
    parent_id BIGINT NULL COMMENT '親カテゴリーID（階層構造用）',
    display_order INT DEFAULT 0 COMMENT '表示順序',
    is_active BOOLEAN DEFAULT TRUE COMMENT '有効フラグ',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_slug (slug),
    INDEX idx_parent_id (parent_id),
    INDEX idx_is_active (is_active),
    INDEX idx_display_order (display_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='カテゴリー';

-- =====================================================
-- 投稿テーブル
-- =====================================================
CREATE TABLE IF NOT EXISTS posts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '投稿者ID',
    category_id BIGINT COMMENT 'カテゴリーID',
    title VARCHAR(255) NOT NULL COMMENT 'タイトル',
    slug VARCHAR(255) NOT NULL COMMENT 'URLスラッグ',
    content TEXT NOT NULL COMMENT '本文',
    excerpt VARCHAR(500) COMMENT '要約',
    status ENUM('draft', 'published', 'archived', 'deleted') NOT NULL DEFAULT 'draft' COMMENT '公開状態',
    published_at TIMESTAMP NULL COMMENT '公開日時',
    view_count INT DEFAULT 0 COMMENT '閲覧数',
    like_count INT DEFAULT 0 COMMENT 'いいね数',
    comment_count INT DEFAULT 0 COMMENT 'コメント数',
    is_featured BOOLEAN DEFAULT FALSE COMMENT '注目記事フラグ',
    meta_title VARCHAR(255) COMMENT 'SEO用タイトル',
    meta_description VARCHAR(500) COMMENT 'SEO用説明',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    UNIQUE INDEX idx_slug (slug),
    INDEX idx_user_id (user_id),
    INDEX idx_category_id (category_id),
    INDEX idx_status (status),
    INDEX idx_published_at (published_at),
    INDEX idx_view_count (view_count),
    INDEX idx_like_count (like_count),
    INDEX idx_is_featured (is_featured),
    INDEX idx_created_at (created_at),
    FULLTEXT INDEX idx_fulltext_search (title, content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投稿';

-- =====================================================
-- タグテーブル
-- =====================================================
CREATE TABLE IF NOT EXISTS tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT 'タグ名',
    slug VARCHAR(50) NOT NULL UNIQUE COMMENT 'URLスラッグ',
    description VARCHAR(255) COMMENT '説明',
    usage_count INT DEFAULT 0 COMMENT '使用回数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_slug (slug),
    INDEX idx_usage_count (usage_count)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='タグ';

-- =====================================================
-- 投稿-タグ中間テーブル（多対多）
-- =====================================================
CREATE TABLE IF NOT EXISTS post_tags (
    post_id BIGINT NOT NULL COMMENT '投稿ID',
    tag_id BIGINT NOT NULL COMMENT 'タグID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    INDEX idx_tag_id (tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投稿タグ関連';

-- =====================================================
-- コメントテーブル
-- =====================================================
CREATE TABLE IF NOT EXISTS comments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    post_id BIGINT NOT NULL COMMENT '投稿ID',
    user_id BIGINT NOT NULL COMMENT 'コメント投稿者ID',
    parent_id BIGINT NULL COMMENT '親コメントID（返信用）',
    content TEXT NOT NULL COMMENT 'コメント内容',
    status ENUM('pending', 'approved', 'rejected', 'spam') NOT NULL DEFAULT 'pending' COMMENT '承認状態',
    like_count INT DEFAULT 0 COMMENT 'いいね数',
    is_edited BOOLEAN DEFAULT FALSE COMMENT '編集済みフラグ',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE,
    INDEX idx_post_id (post_id),
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='コメント';

-- =====================================================
-- フォロー関係テーブル（多対多）
-- =====================================================
CREATE TABLE IF NOT EXISTS user_follows (
    follower_id BIGINT NOT NULL COMMENT 'フォローする側のユーザーID',
    following_id BIGINT NOT NULL COMMENT 'フォローされる側のユーザーID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_following_id (following_id),
    CHECK (follower_id != following_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザーフォロー関係';

-- =====================================================
-- いいねテーブル
-- =====================================================
CREATE TABLE IF NOT EXISTS likes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT 'いいねしたユーザーID',
    likeable_type ENUM('post', 'comment') NOT NULL COMMENT 'いいね対象種別',
    likeable_id BIGINT NOT NULL COMMENT 'いいね対象ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE INDEX idx_user_likeable (user_id, likeable_type, likeable_id),
    INDEX idx_likeable (likeable_type, likeable_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='いいね';

-- =====================================================
-- 通知テーブル
-- =====================================================
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '通知先ユーザーID',
    type ENUM('follow', 'like', 'comment', 'mention', 'system') NOT NULL COMMENT '通知タイプ',
    title VARCHAR(255) NOT NULL COMMENT 'タイトル',
    message TEXT NOT NULL COMMENT 'メッセージ',
    link_url VARCHAR(500) COMMENT 'リンクURL',
    is_read BOOLEAN DEFAULT FALSE COMMENT '既読フラグ',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP NULL COMMENT '既読日時',
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知';

-- =====================================================
-- サンプルデータ投入
-- =====================================================

-- ユーザーデータ
INSERT INTO users (username, email, password_hash, status, email_verified, last_login_at) VALUES
    ('sakura', 'sakura@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', TRUE, NOW() - INTERVAL 1 HOUR),
    ('takeshi', 'takeshi@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', TRUE, NOW() - INTERVAL 2 DAY),
    ('yuki', 'yuki@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', TRUE, NOW() - INTERVAL 5 DAY),
    ('haruto', 'haruto@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'inactive', FALSE, NULL)
ON DUPLICATE KEY UPDATE username=VALUES(username);

-- プロフィールデータ
INSERT INTO user_profiles (user_id, first_name, last_name, display_name, bio, birth_date, gender, country_code, timezone) VALUES
    (1, '桜', '田中', 'さくら', 'クリーンアーキテクチャとフルスタック開発に情熱を注いでいます。', '1995-05-15', 'female', 'JP', 'Asia/Tokyo'),
    (2, '健', '佐藤', 'たけし', 'DevOpsエンジニアとして自動化とクラウドインフラを愛しています。', '1988-08-22', 'male', 'JP', 'Asia/Tokyo'),
    (3, '雪', '鈴木', 'ゆき', '機械学習とAIを探求するデータサイエンティストです。', '1982-03-10', 'female', 'JP', 'Asia/Tokyo'),
    (4, '陽翔', '高橋', 'はると', 'マイクロサービスを専門とするバックエンドエンジニアです。', '1992-11-30', 'male', 'JP', 'Asia/Tokyo')
ON DUPLICATE KEY UPDATE user_id=VALUES(user_id);

-- カテゴリーデータ
INSERT INTO categories (name, slug, description, parent_id, display_order) VALUES
    ('テクノロジー', 'technology', 'テクノロジーとイノベーションに関するすべて', NULL, 1),
    ('プログラミング', 'programming', 'プログラミング言語と技術', 1, 1),
    ('DevOps', 'devops', 'DevOpsの実践とツール', 1, 2),
    ('ライフスタイル', 'lifestyle', 'ライフスタイルと日常生活', NULL, 2),
    ('旅行', 'travel', '旅行体験とヒント', 4, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- タグデータ
INSERT INTO tags (name, slug, description, usage_count) VALUES
    ('Go言語', 'go', 'Go プログラミング言語', 5),
    ('Docker', 'docker', 'Docker コンテナ化', 8),
    ('Kubernetes', 'kubernetes', 'Kubernetes オーケストレーション', 3),
    ('MySQL', 'mysql', 'MySQL データベース', 7),
    ('Redis', 'redis', 'Redis キャッシング', 4),
    ('クリーンアーキテクチャ', 'clean-architecture', 'クリーンアーキテクチャパターン', 6),
    ('マイクロサービス', 'microservices', 'マイクロサービスアーキテクチャ', 5),
    ('テスト', 'testing', 'ソフトウェアテスト', 4)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 投稿データ
INSERT INTO posts (user_id, category_id, title, slug, content, excerpt, status, published_at, view_count, like_count, comment_count, is_featured) VALUES
    (1, 2, 'Go言語でクリーンアーキテクチャを始めよう', 'clean-architecture-go', 
     'クリーンアーキテクチャは、設計要素をリング状のレベルに分離するソフトウェア設計哲学です。この記事では、Goでの実装方法を探ります...', 
     'Goアプリケーションでクリーンアーキテクチャパターンを実装する方法を学びましょう', 
     'published', NOW() - INTERVAL 10 DAY, 1250, 45, 12, TRUE),
    
    (2, 3, 'DockerとKubernetesのベストプラクティス', 'docker-kubernetes-best-practices',
     'コンテナ化はソフトウェアデプロイメントに革命をもたらしました。この包括的なガイドでは、DockerとKubernetesのベストプラクティスを探ります...', 
     'DockerとKubernetesのデプロイメント戦略に関する包括的なガイド',
     'published', NOW() - INTERVAL 5 DAY, 890, 32, 8, TRUE),
    
    (1, 2, 'MySQLパフォーマンス最適化のヒント', 'mysql-performance-optimization',
     'データベースのパフォーマンスはアプリケーションの成功に不可欠です。ここでは、MySQLクエリを最適化するための実証済みのテクニックを紹介します...', 
     'より良いパフォーマンスのための必須のMySQL最適化テクニック',
     'published', NOW() - INTERVAL 3 DAY, 654, 28, 15, FALSE),
    
    (3, 2, '高トラフィックアプリケーションのためのRedisキャッシング戦略', 'redis-caching-strategies',
     'キャッシングはWebアプリケーションのスケーリングに不可欠です。この記事では、さまざまなRedisキャッシングパターンについて説明します...', 
     'スケーラブルなアプリケーションのための効果的なRedisキャッシングパターンを学びましょう',
     'published', NOW() - INTERVAL 1 DAY, 423, 19, 7, FALSE),
    
    (2, 3, 'Goでマイクロサービスを構築する', 'building-microservices-go',
     'マイクロサービスアーキテクチャは柔軟性とスケーラビリティを提供します。Goを使用してマイクロサービスを構築しましょう...', 
     'Goでマイクロサービスを作成するためのステップバイステップガイド',
     'draft', NULL, 0, 0, 0, FALSE)
ON DUPLICATE KEY UPDATE slug=VALUES(slug);

-- 投稿-タグ関連データ
INSERT INTO post_tags (post_id, tag_id) VALUES
    (1, 1), (1, 6),  -- Clean Architecture: Go, CleanArchitecture
    (2, 2), (2, 3),  -- Docker K8s: Docker, Kubernetes
    (3, 4),          -- MySQL: MySQL
    (4, 5), (4, 1),  -- Redis: Redis, Go
    (5, 1), (5, 7)   -- Microservices: Go, Microservices
ON DUPLICATE KEY UPDATE post_id=VALUES(post_id);

-- コメントデータ
INSERT INTO comments (post_id, user_id, parent_id, content, status, like_count) VALUES
    (1, 2, NULL, '素晴らしい記事です！クリーンアーキテクチャの実践的なガイドを探していました。', 'approved', 5),
    (1, 3, 1, '同感です！例がとても分かりやすいですね。', 'approved', 2),
    (1, 4, NULL, 'インフラストラクチャ層についてもっと例を提供していただけますか？', 'approved', 1),
    (2, 1, NULL, 'Kubernetesのベストプラクティスの素晴らしい解説ですね！', 'approved', 8),
    (2, 3, NULL, 'セキュリティのセクションが特に役立ちました。', 'approved', 3),
    (3, 2, NULL, 'これらの最適化のヒントが私のプロジェクトを救いました。ありがとうございます！', 'approved', 6),
    (4, 1, NULL, 'とてもタイムリーな記事です。ちょうどRedisキャッシングを実装したところでした。', 'approved', 4)
ON DUPLICATE KEY UPDATE post_id=VALUES(post_id);

-- フォロー関係データ
INSERT INTO user_follows (follower_id, following_id) VALUES
    (1, 2), (1, 3),  -- さくらがたけしとゆきをフォロー
    (2, 1), (2, 3),  -- たけしがさくらとゆきをフォロー
    (3, 1), (3, 2),  -- ゆきがさくらとたけしをフォロー
    (4, 1), (4, 2), (4, 3)  -- はるとが全員をフォロー
ON DUPLICATE KEY UPDATE follower_id=VALUES(follower_id);

-- いいねデータ
INSERT INTO likes (user_id, likeable_type, likeable_id) VALUES
    (2, 'post', 1), (3, 'post', 1), (4, 'post', 1),  -- 投稿1へのいいね
    (1, 'post', 2), (3, 'post', 2),  -- 投稿2へのいいね
    (1, 'comment', 1), (3, 'comment', 1)  -- コメントへのいいね
ON DUPLICATE KEY UPDATE user_id=VALUES(user_id);

-- 通知データ
INSERT INTO notifications (user_id, type, title, message, link_url, is_read) VALUES
    (1, 'follow', '新しいフォロワー', 'たけしさんがあなたをフォローしました', '/users/takeshi', FALSE),
    (1, 'like', '新しいいいね', 'ゆきさんがあなたの投稿にいいねしました', '/posts/1', TRUE),
    (1, 'comment', '新しいコメント', 'はるとさんがあなたの投稿にコメントしました', '/posts/1#comment-3', FALSE),
    (2, 'like', '新しいいいね', 'さくらさんがあなたの投稿にいいねしました', '/posts/2', FALSE)
ON DUPLICATE KEY UPDATE user_id=VALUES(user_id);
