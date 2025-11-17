package handlers

import "github/jornal-cidadao-jc/internal/storage"

type Handler struct {
    Storage    *storage.Storage
    ChargesDir string
    PostsDir string
}

func NewHandler(s *storage.Storage, chargesDir string, postsDir string) *Handler {
    return &Handler{
        Storage:    s,
        ChargesDir: chargesDir,
        PostsDir: postsDir,
    }
}
