package handlers

import (
	"database/sql"
	"fmt"
	"github/jornal-cidadao-jc/internal/model"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateStatusPayload struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.Storage.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro processando lista de usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) UploadCharge(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "O campo de título é obrigatório.",
		})
		return
	}

	file, err := c.FormFile("charge_file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Nenhum arquivo foi enviado. Por favor, selecione um arquivo.",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Formato de arquivo inválido. Apenas imagens (.png, .jpg, .jpeg, .gif) são permitidas.",
		})
		return
	}

	uniqueFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	destinationPath := filepath.Join(h.ChargesDir, uniqueFilename)

	if err := c.SaveUploadedFile(file, destinationPath); err != nil {
		log.Println("Erro ao salvar o arquivo:", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Não foi possível salvar o arquivo no servidor. Tente novamente mais tarde.",
		})
		return
	}

	if err := h.Storage.CreateCharge(title, uniqueFilename); err != nil {
		log.Println("Erro ao salvar metadados no banco:", err)
		os.Remove(destinationPath)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "O arquivo foi salvo, mas houve um erro ao registrar as informações no banco de dados.",
		})
		return
	}

	c.HTML(http.StatusOK, "adicionar_charge.tmpl", gin.H{
		"message": "Charge '" + title + "' enviada com sucesso!",
	})
}

func (h *Handler) DeleteCharge(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Erro ao obter ID da charge: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return 
	}

	filename, err := h.Storage.DeleteCharge(id)
	if err != nil {
		log.Println("Erro ao deletar charge do banco de dados: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar charge do banco de dados"})
		return 
	}

	destinationPath := filepath.Join(h.ChargesDir, filename)
	err = os.Remove(destinationPath)
	if err != nil {
		log.Println("Erro ao deletar arquivo físico da charge: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar o arquivo físico, mas o registro no banco foi removido"})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"message": "Charge '" + filename + "' deletada com sucesso!"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Erro ao obter ID do usuário: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return 
	}

	user, err := h.Storage.GetUserByID(id)
	if err != nil {
		log.Println("Erro ao obter usuário do banco de dados: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter usuário do banco de dados"})
		return 
	}

	err = h.Storage.DeleteUser(id)
	if err != nil {
		log.Println("Erro ao deletar usuário do banco de dados: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar usuário do banco de dados"})
		return 
	}


	c.JSON(http.StatusOK, gin.H{"message": "Usuário '" + user.Username + "' deletada com sucesso!"})
}

func (h *Handler) UploadArticle(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	body := c.PostForm("body")
	
	if title == "" || author == "" || body == "" {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "Título, autor e corpo da matéria são obrigatórios."})
		return
	}

	pollQuestion := c.PostForm("poll_question")
	pollOptions := c.PostFormArray("poll_options")

	log.Println("pollQuestion : ", pollQuestion )
	log.Println("pollOptions: ", pollOptions)

	err := h.Storage.CreateArticleWithPoll(title, author, body, pollQuestion, pollOptions)
	if err != nil {
		log.Println("Erro salvando artigo com enquete no banco de dados: ", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Erro ao salvar o artigo."})
		return
	}

	c.HTML(http.StatusOK, "escrever_materia.tmpl", gin.H{
		"message": "Matéria '" + title + "' criada com sucesso!",
	})
}

func (h *Handler) UpdateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "Id da matéria inválido"})
		return
	}
	
	title := c.PostForm("title")
	author := c.PostForm("author")
	body := c.PostForm("body")

	if title == "" || author == "" || body == "" {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": "Título, autor e corpo da matéria são obrigatórios."})
		return
	}

	pollQuestion := c.PostForm("poll_question")
	pollOptions := c.PostFormArray("poll_options")

	err = h.Storage.UpdateArticleWithPoll(id, title, author, body, pollQuestion, pollOptions)
	if err != nil {
		log.Println("Erro atualizando artigo com enquete no banco de dados: ", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Erro ao atualizar o artigo."})
		return
	}

	c.HTML(http.StatusOK, "atualizar_materia.tmpl", gin.H{
		"message": "Matéria '" + title + "' atualizada com sucesso!",
	})
}

func (h *Handler) UpdatePostStatus(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post inválido."})
		return
	}

	newStatus := model.ToPostStatus(c.Param("status"))

	err = h.Storage.UpdatePostStatus(postID, newStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post não encontrado."})
			return
		}
		log.Println("Erro ao atualizar o status do post:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno ao atualizar o status."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status do post atualizado com sucesso!"})
}
