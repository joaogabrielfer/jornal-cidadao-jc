package handlers

import(
	"github.com/gin-gonic/gin"
	"math/rand/v2"
	"strconv"
	"path/filepath"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"log"
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


func (h *Handler) GetArticles(c *gin.Context){
	articles, err := h.Storage.GetArticles()
	if err != nil{
		log.Println("erro buscando artigos no banco de dados: ", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error": "Erro buscando artigos no banco de dados",
		})
	}

	c.JSON(http.StatusOK, articles)
}
func JornalCidadaoTest(c *gin.Context) {
    c.HTML(http.StatusOK, "ultimas.tmpl", gin.H{
        "AnoAtual": 2025,
        "Artigos":  nil,
    })
}
func (h *Handler) GetArticleByID(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		log.Println("erro convertendo id para inteiro: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Argumento de ID inválido",
		})
		return
	}

	article, err := h.Storage.GetArticleByID(id)
	if err != nil{
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

func (h *Handler) UpdateVoteCount(c *gin.Context){
	optionID, err :=strconv.Atoi(c.Param("id"))
	if err != nil{
		log.Println("Erro ao obter ID da opção da enquete: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID inválido"})
		return 
	}

	err = h.Storage.VotePoll(optionID)
	if err != nil{
		log.Println("Erro registrando voto no banco de dados: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro registrando voto no banco de dados"})
		return 
	}
	c.JSON(http.StatusOK, gin.H{"message": "Voto registrado com sucesso"})
}
