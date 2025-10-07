package handlers

import (
	"fmt"
	"log"
	"database/sql"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github/jornal-cidadao-jc/internal/storage"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Storage    *storage.Storage
	ChargesDir string
}

func NewHandler(s *storage.Storage, chargesDir string) *Handler {
	return &Handler{
		Storage:    s,
		ChargesDir: chargesDir,
	}
}

func (h *Handler) GetIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
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

func (h *Handler) CreateUser(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if c.PostForm("password") != c.PostForm("password-confirm") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A senha na confirmaçao esta diferente."})
		return
	}
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
	c.Redirect(http.StatusSeeOther, "/cadastro")
}

func (h *Handler) GetSignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "cadastro.tmpl", nil)
}

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.Storage.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro processando lista de usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
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

func (h *Handler) GetAdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.tmpl", nil)
}

func (h *Handler) GetUsersAdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_users.tmpl", nil)
}

func (h *Handler) GetUploadChargePage(c *gin.Context) {
	c.HTML(http.StatusOK, "adicionar_charge.tmpl", nil)
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

func (h *Handler) GetDeleteChargePage(c *gin.Context) {
	c.HTML(http.StatusOK, "deletar_charge.tmpl", nil)
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
