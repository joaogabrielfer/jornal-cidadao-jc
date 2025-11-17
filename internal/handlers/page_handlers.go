package handlers

import (
	"database/sql"
	"html/template"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github/jornal-cidadao-jc/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetIndexPage(c *gin.Context)      { c.HTML(http.StatusOK, "index.tmpl", nil) }
func (h *Handler) GetLoginPage(c *gin.Context)      { c.HTML(http.StatusOK, "login.tmpl", nil) }
func (h *Handler) GetSignupPage(c *gin.Context)     { c.HTML(http.StatusOK, "cadastro.tmpl", nil) }
func (h *Handler) GetAdminPage(c *gin.Context)      { c.HTML(http.StatusOK, "admin.tmpl", nil) }
func (h *Handler) GetUsersAdminPage(c *gin.Context) { c.HTML(http.StatusOK, "admin_users.tmpl", nil) }

func (h *Handler) GetUploadChargePage(c *gin.Context) {
	c.HTML(http.StatusOK, "adicionar_charge.tmpl", nil)
}
func (h *Handler) GetDeleteChargePage(c *gin.Context) {
	c.HTML(http.StatusOK, "deletar_charge.tmpl", nil)
}
func (h *Handler) GetUploadArticlePage(c *gin.Context) {
	c.HTML(http.StatusOK, "escrever_materia.tmpl", nil)
}
func (h *Handler) GetUpdateArticlePage(c *gin.Context) {
	c.HTML(http.StatusOK, "atualizar_materia.tmpl", nil)
}

func (h *Handler) GetArticlesPage(c *gin.Context) {
	articles, err := h.Storage.GetArticles()
	if err != nil {
		log.Println("Erro ao buscar matérias para a página de admin:", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Não foi possível carregar as matérias."})
		return
	}

	c.HTML(http.StatusOK, "materias.tmpl", gin.H{
		"Articles": articles,
	})
}

func (h *Handler) ShowJornalCidadaoDashboard(c *gin.Context) {
    authorID := 1 
    posts, err := h.Storage.GetPostsByAuthorID(authorID)
    if err != nil && err != sql.ErrNoRows {
        log.Println("Erro ao buscar posts do autor:", err)
        c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"error": "Não foi possível carregar suas notícias."})
        return
    }
    
    if posts == nil {
        posts = []model.Post{}
    }

    c.HTML(http.StatusOK, "jc.tmpl", gin.H{
        "Posts": posts, 
    })
}

func (h *Handler) GetNoIdChargePage(c *gin.Context) {
	charges, err := h.Storage.GetAllCharges()
	if err != nil {
		log.Println("Erro obtendo charges: ", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Não foi possível obter charges",
		})
		return
	}
	highest_id := -1
	for _, charge := range charges {
		if charge.ID > highest_id {
			highest_id = charge.ID
		}
	}
	if highest_id == -1 {
		log.Println("Erro obtendo ID de maior charge")
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Não foi possível obter id da charge",
		})
	} else {
		highest_id_str := strconv.Itoa(highest_id)
		c.Redirect(http.StatusSeeOther, "/charge/"+highest_id_str)
	}
}

func (h *Handler) GetChargePage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Erro convertendo ID para integer", err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "O ID fornecido na URL é inválido.",
		})
		return
	}
	charge, err := h.Storage.GetChargeByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
				"error": "A charge com este ID não foi encontrada.",
			})
			return
		}
		log.Println("Erro obtendo charge do DB por ID: ", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Ocorreu um erro inesperado ao buscar a charge.",
		})
		return
	}

	charge.URL = filepath.Join("/static/images/charges", charge.Filename)
	c.HTML(http.StatusOK, "vizualizar_charge.tmpl", gin.H{
		"charge": charge,
	})
}
