package postgres

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	pb "post-servic/genproto/post"
	"post-servic/storage"
)

type BasicAdditional struct {
	db *sqlx.DB
}

func NewBasicAdditional(db *sqlx.DB) storage.BasicAdditional {
	return &BasicAdditional{db}
}

func (b *BasicAdditional) GetUserRecommendation(in *pb.Username) (*pb.PostListResponse, error) {
	res := &pb.PostListResponse{
		Post: []*pb.PostResponse{},
	}

	var postID []string
	err := b.db.Select(&postID, "SELECT post_id FROM likes WHERE user_id=$1 ORDER BY created_at DESC LIMIT 10", in.Username)
	if errors.Is(err, sql.ErrNoRows) {
		err = b.db.Select(&res.Post, "SELECT id, user_id, country, location, title, hashtag, content, image_url, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT 20")
		return res, err
	}
	if err != nil {
		log.Println("Error in getting Post ID:", err)
		return nil, err
	}

	// Получаем user_id из таблицы постов
	var userId []string
	q := "SELECT user_id FROM posts WHERE id IN (?)"
	query, args, err := sqlx.In(q, postID)
	if err != nil {
		log.Println("Error preparing request for user IDs:", err)
		return nil, err
	}

	query = b.db.Rebind(query)
	err = b.db.Select(&userId, query, args...)
	if err != nil {
		log.Println("Error in getting User ID:", err)
		return nil, err
	}

	// Получаем национальности из таблицы постов
	var nationality []string
	q = "SELECT country FROM posts WHERE id IN (?)"
	query, args, err = sqlx.In(q, postID)
	if err != nil {
		log.Println("Error preparing request for counties:", err)
		return nil, err
	}

	query = b.db.Rebind(query)
	err = b.db.Select(&nationality, query, args...)
	if err != nil {
		log.Println("Error in getting Countries:", err)
		return nil, err
	}

	// Получаем хэштеги из таблицы постов
	var hashtag []string
	q = "SELECT hashtag FROM posts WHERE id IN (?)"
	query, args, err = sqlx.In(q, postID)
	if err != nil {
		log.Println("Error preparing request for hashtag:", err)
		return nil, err
	}

	query = b.db.Rebind(query)
	err = b.db.Select(&hashtag, query, args...)
	if err != nil {
		log.Println("Error in getting hashtag:", err)
		return nil, err
	}

	// Получаем рекомендованные посты для пользователя
	q = `SELECT id, user_id, country, location, title, hashtag, content, image_url, created_at, updated_at 
         FROM posts 
         WHERE country IN (?) AND user_id IN (?) AND hashtag IN (?) 
         ORDER BY created_at DESC
         LIMIT 20`
	query, args, err = sqlx.In(q, nationality, userId, hashtag)
	if err != nil {
		log.Println("Error preparing request for getting Posts:", err)
		return nil, err
	}

	query = b.db.Rebind(query)
	rows, err := b.db.Query(query, args...)
	if err != nil {
		log.Println("Error in getting posts:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post pb.PostResponse
		err = rows.Scan(&post.Id, &post.UserId, &post.Country, &post.Location, &post.Title, &post.Hashtag, &post.Content, &post.ImageUrl, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			log.Println("Error in Scanning posts:", err)
			return nil, err
		}

		res.Post = append(res.Post, &post)
	}

	if err = rows.Err(); err != nil {
		log.Println("Ошибка при переборе строк:", err)
		return nil, err
	}

	return res, nil
}

func (b *BasicAdditional) GetPostsByUsername(in *pb.Username) (*pb.PostListResponse, error) {
	var res *pb.PostListResponse

	err := b.db.Select(&res, `SELECT id, user_id, nationality, location, title, hashtag, content, image_url,
       created_at, updated_at FROM posts WHERE user_id=$1 order by created_at desc`, in.Username)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BasicAdditional) GetTrendsPost(in *pb.Void) (*pb.PostListResponse, error) {
	return nil, nil
}

func (b *BasicAdditional) SearchPost(in *pb.Search) (*pb.PostListResponse, error) {
	var res *pb.PostListResponse

	if in.Action == "" {
		return res, nil
	}

	in.Action = "%" + in.Action + "%"

	query := `SELECT id, user_id, country, location, title, hashtag, content, image_url,
       created_at, updated_at FROM posts WHERE title ilike $1 or description ilike $2 or hashtag $3`

	err := b.db.Select(&res, query, in.Action, in.Action, in.Action)
	if err != nil {
		return nil, err
	}

	return res, nil
}
