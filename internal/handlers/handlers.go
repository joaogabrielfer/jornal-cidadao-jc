package handlers

import (
	"log"
	"math/rand/v2"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github/jornal-cidadao-jc/internal/model"
	"github/jornal-cidadao-jc/internal/storage"
	"github/jornal-cidadao-jc/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Storage      *storage.Storage
	ChargesDir   string
}

func New_handler(s *storage.Storage, chargesDir string) *Handler {
	return &Handler{
		Storage:      s,
		ChargesDir:   chargesDir,
	}
}

func (h *Handler) Get_index_page(c *gin.Context){
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func (h *Handler) Get_charge_page(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		log.Println("Erro convertendo ID para integer", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de charge inválido"})
		return
	}
	charges, err := utils.Get_charges_object(h.ChargesDir)
	if err != nil{
		log.Println("Erro obtendo informações das charges: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro obtendo informações das charges"})
		return
	}

	var chosen_charge model.Charge
	for _, charge := range charges{
		if charge.ID == id{
			chosen_charge = charge
		}
	}

	c.HTML(http.StatusOK, "vizualizar_charge.tmpl", gin.H{
		"charge":chosen_charge,
	})
}

func (h *Handler) Create_user(c *gin.Context) {
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

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro fazendo hash da senha: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falhou em criar conta"})
		return
	}

	err = h.Storage.Create_user(username, email, string(hashed_password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome ou email ja podem estar em uso"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/cadastro")
}

func (h *Handler) Get_signup_page(c *gin.Context) {
	c.HTML(http.StatusOK, "cadastro.tmpl", nil)
}

func (h *Handler) Get_users(c *gin.Context) {
	users, err := h.Storage.Get_users()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro processando lista de usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) Get_charges_list(c *gin.Context) {
	charges, err := utils.Get_charges_object(h.ChargesDir)
	if err != nil{
		log.Println("Erro obtendo informações das charges: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro obtendo informações das charges"})
		return
	}
	if len(charges) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma charge encontrada"})
		return
	}

	sort.Slice(charges, func(i, j int) bool {
		return time.Time(charges[i].Date).After(time.Time(charges[j].Date))
	})

	c.JSON(http.StatusOK, charges)
}

func (h *Handler) Get_random_charge(c *gin.Context){
	charges, err := utils.Get_charges_object(h.ChargesDir)
	if err != nil{
		log.Println("Erro obtendo informações das charges: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro obtendo informações das charges"})
		return
	}
	if len(charges) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma charge encontrada"})
		return
	}

	var chosen_charge model.Charge
	random_index := rand.N(len(charges)) + 1
	for _, charge := range charges{
		if charge.ID == random_index{
			chosen_charge = charge
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"charge": chosen_charge,
	})
}
