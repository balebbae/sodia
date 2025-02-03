package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/balebbae/sodia/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title string  `json:"title"`
	Content string `json:"content"`
	Tags []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
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
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return 
	}

	if err = writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return 
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return 
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return 
	}
	if err = writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return 
	} 
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the post ID from the URL parameters.
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Decode the JSON payload into the payload struct.
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Retrieve the existing post.
	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "Post not found")
		} else {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch post")
		}
		return
	}

	// Update the post fields.
	post.Title = payload.Title
	post.Content = payload.Content
	post.Tags = payload.Tags
	// TODO: Update the UserID after implementing authentication.
	post.UserID = 1

	// Save the updated post.
	if err := app.store.Posts.Update(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	// Send the updated post as JSON.
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
