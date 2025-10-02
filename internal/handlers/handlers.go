package handlers

import (
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"path"

	"github/jornal-cidadao-jc/internal/model"
	"github/jornal-cidadao-jc/internal/storage"

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
	files, err := os.ReadDir(h.ChargesDir)
	if err != nil {
		log.Println("Erro lendo diretorio de charges: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível ler diretório de charges."})
		return
	}

	var charges []model.ChargesInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		file_path := filepath.Join(h.ChargesDir, file.Name())
		file_info, err := os.Stat(file_path)
		if err != nil {
			log.Println("Erro obtendo informação da charge: ", file.Name(), err)
			continue
		}
		charges = append(charges, model.ChargesInfo{
			Filename: file.Name(),
			ModTime:  file_info.ModTime(),
		})
	}

	if len(charges) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhuma charge encontrada"})
		return
	}

	sort.Slice(charges, func(i, j int) bool {
		return charges[i].ModTime.After(charges[j].ModTime)
	})

	var response_data []model.ChargeResponse
	for _, charge := range charges {
		response_data = append(response_data, model.ChargeResponse{
			Filename: charge.Filename,
			Date:     charge.ModTime.Format("02-01-2006 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, response_data)
}

func (h *Handler) Get_index_page(c *gin.Context){
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func (h *Handler) Get_random_charge(c *gin.Context){
	dir_entries, err := os.ReadDir(h.ChargesDir)
	if err != nil {
		log.Println("Erro lendo diretorio de charges: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível ler diretório de charges."})
		return
	}

	var files []os.DirEntry
	for _, entry := range dir_entries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}

	if len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhuma charge (arquivo) encontrada"})
		return
	}
	var chosen_charge model.ChargesInfo

	random_index := rand.N(len(files))
	file_info, err := files[random_index].Info()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível obter informações da charge"})
	}

	chosen_charge.Filename = files[random_index].Name()
	chosen_charge.URL = path.Join("/static", "images", "charges", files[random_index].Name())
	chosen_charge.ModTime = file_info.ModTime()
	chosen_charge.Title = ""

	c.JSON(http.StatusOK, gin.H{
		"charge": chosen_charge,
	})
}
