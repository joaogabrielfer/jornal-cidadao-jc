package handlers

import (
	"fmt"
	"github/jornal-cidadao-jc/internal/model"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUser(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if c.PostForm("password") != c.PostForm("password-confirm") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A senha na confirmaçao esta diferente."})
		return
	}

	log.Printf("username = %s, email = %s, password = %s, password-confirm = %s\n", username, email, password, c.PostForm("password-confirm"))
	if username == "" || email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todos os campos sao requeridos"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro fazendo hash da senha: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falhou em criar conta"})
		return
	}

	err = h.Storage.CreateUser(username, email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome ou email ja podem estar em uso"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/login?signup=success")
}

func (h *Handler) GetChargesList(c *gin.Context) {
	charges, err := h.Storage.GetAllCharges()
	if err != nil {
		log.Println("Erro ao buscar charges no banco de dados:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter lista de charges."})
		return
	}

	for i := range charges {
		charges[i].URL = filepath.Join("/static/images/charges", charges[i].Filename)
	}
	c.JSON(http.StatusOK, charges)
}

func (h *Handler) GetRandomCharge(c *gin.Context) {
	charges, err := h.Storage.GetAllCharges()
	if err != nil {
		log.Println("Erro obtendo informações das charges: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro obtendo informações das charges"})
		return
	}
	if len(charges) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma charge encontrada"})
		return
	}

	randomIndex := rand.N(len(charges))
	chosenCharge := charges[randomIndex]
	chosenCharge.URL = filepath.Join("/static/images/charges", chosenCharge.Filename)

	c.JSON(http.StatusOK, gin.H{
		"charge": chosenCharge,
	})
}

func (h *Handler) GetArticles(c *gin.Context) {
	articles, err := h.Storage.GetArticles()
	if err != nil {
		log.Println("erro buscando artigos no banco de dados: ", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Erro buscando artigos no banco de dados",
		})
	}

	c.JSON(http.StatusOK, articles)
}

func (h *Handler) GetArticleByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("erro convertendo id para inteiro: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Argumento de ID inválido",
		})
		return
	}

	article, err := h.Storage.GetArticleByID(id)
	if err != nil {
		log.Println("erro obtendo artigo", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro obtendo artigo",
		})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *Handler) DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Erro ao obter ID da matéria: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	article, err := h.Storage.GetArticleByID(id)
	if err != nil {
		log.Println("Erro ao obter matéria do banco de dados: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter matéria do banco de dados"})
		return
	}

	err = h.Storage.DeleteArticle(id)
	if err != nil {
		log.Println("Erro ao deletar matéria do banco de dados: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar matéria do banco de dados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Matéria '" + article.Title + "' deletada com sucesso!"})
}

func (h *Handler) UpdateVoteCount(c *gin.Context) {
	optionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Erro ao obter ID da opção da enquete: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID inválido"})
		return
	}

	err = h.Storage.VotePoll(optionID)
	if err != nil {
		log.Println("Erro registrando voto no banco de dados: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro registrando voto no banco de dados"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Voto registrado com sucesso"})
}

func (h *Handler) UploadPost(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")

	if title == "" || description == "" {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Os campos de título e descrição são obrigatórios.",
		})
		return
	}

	//placeholder - como ainda não tem login
	authorID := 1

	file, err := c.FormFile("media_file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Nenhum arquivo de mídia foi enviado.",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" && ext != ".mp4" {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Formato de arquivo inválido. Apenas imagens (png, jpg, jpeg, gif) e vídeos (mp4) são permitidos.",
		})
		return
	}

	uniqueFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	destinationPath := filepath.Join(h.PostsDir, uniqueFilename)

	if err := c.SaveUploadedFile(file, destinationPath); err != nil {
		log.Println("Erro ao salvar o arquivo de mídia:", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Não foi possível salvar o arquivo no servidor.",
		})
		return
	}

	mediaURL := filepath.ToSlash(filepath.Join("/static", "media", "posts", uniqueFilename))

	if err := h.Storage.CreatePost(title, description, mediaURL, authorID, model.StatusEmAnalise); err != nil {
		log.Println("Erro ao salvar metadados do post no banco:", err)
		os.Remove(destinationPath)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "O arquivo foi salvo, mas houve um erro ao registrar o post no banco de dados.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sucess",
		"message": "Notícia enviada para a moderação.",
	})
}

func (h *Handler) GetApprovedPosts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	posts, metadata, err := h.Storage.GetApprovedPostsPaginated(page, limit)
	if err != nil {
		log.Println("Erro ao buscar posts paginados no banco de dados:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter lista de posts."})
		return
	}

	response := model.PaginatedPosts{
		Posts:    posts,
		Metadata: metadata,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("erro convertendo id para inteiro: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Argumento de ID inválido",
		})
		return
	}

	post, err := h.Storage.GetPostByID(id)
	if err != nil {
		log.Println("erro obtendo post", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro obtendo post",
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetPostByAuthorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("erro convertendo id para inteiro: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Argumento de ID inválido",
		})
		return
	}

	posts, err := h.Storage.GetPostsByAuthorID(id)
	if err != nil {
		log.Println("erro obtendo post", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro obtendo posts do usuário",
		})
		return
	}

	if posts == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetChargeDoDia(c *gin.Context) {
	charges, err := h.Storage.GetAllCharges()
	if err != nil {
		log.Println("Erro ao buscar charges:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter charge do dia."})
		return
	}

	if len(charges) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma charge encontrada."})
		return
	}

	chargeDoDia := charges[0] 
	chargeDoDia.URL = filepath.Join("/static/images/charges", chargeDoDia.Filename)

	c.JSON(http.StatusOK, chargeDoDia)
}

func (h *Handler) ReportPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post inválido"})
		return
	}

	var payload struct {
		Reason string `json:"reason" form:"reason"`
	}

	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido"})
		return
	}

	if payload.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O motivo da denúncia é obrigatório"})
		return
	}

	// Verifica se o post existe antes de denunciar
	_, err = h.Storage.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post não encontrado"})
		return
	}

	err = h.Storage.CreatePostReport(postID, payload.Reason)
	if err != nil {
		log.Println("Erro ao criar denúncia:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar denúncia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Denúncia enviada com sucesso. Obrigado pela colaboração."})
}
