package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/balebbae/sodia/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string
const postCtx postKey = "post"

type CreatePostPayload struct {
	Title string  `json:"title" validate:"required,max=100"` // validator 
	Content string `json:"content" validate:"required,max=1000"`
	Tags []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return 
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content == "" {
		app.badRequestResponse(w,r, fmt.Errorf("content is required"))
		return
	}

	post := &store.Post{
		Title: payload.Title,
		Content: payload.Content,
		Tags: payload.Tags,
		// TODO: Change after auth
		UserID: 1,
	}

	ctx := r.Context()

	err := app.store.Posts.Create(ctx, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return 
	}

	if err = writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return 
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w,r, err)
		return 	
	}

	// RICH DATA 
	post.Comments = comments

	if err = writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return 
	} 
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the post ID from the URL parameters.
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	err = app.store.Posts.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=100"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = Validate.Struct(payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	err = app.store.Posts.Update(r.Context(), post)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	err = writeJSON(w, http.StatusOK, post)
	if err != nil {
		app.internalServerError(w, r, err) 
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return 
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return 
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}