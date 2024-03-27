package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Address struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep          string `json:"cep"`
	State        string `json:"uf"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
	Complemento  string `json:"complemento"`
	Ibge         string `json:"ibge"`
	Gia          string `json:"gia"`
	DDD          string `json:"ddd"`
	Siafi        string `json:"siafi"`
}

func main() {
	c1 := make(chan Address)
	c2 := make(chan Address)

	var search string
	for _, i := range os.Args[1:] {
		search = search + i
	}

	// VALIDAÇÃO DO CEP DIGITADO
	// REMOVER O HÍFEN, SE HOUVER
	search = strings.Replace(search, "-", "", 1)
	// SE É APENAS NÚMEROS
	match, _ := regexp.MatchString("[0-9]", search)
	if !match {
		fmt.Printf("Informe um CEP Válido. Exemplo: go run main.go 13330-250\n")
		fmt.Printf("Consulte o README.md\n")
		return
	}
	// SE POSSUÍ O NUMERO CORRETO DE CARACTERES
	if search == "" || len(search) != 8 {
		fmt.Printf("Informe um CEP. Exemplo: go run main.go 13330-250\n")
		fmt.Printf("Consulte o README.md\n")
		return
	}

	go SearchInViaCep(search, c1)
	go SearchInBrasilApi(search, c2)

	select {
	case msg := <-c1:
		fmt.Printf("Resultado encontrado: %v - %v (Fonte: ViaCep)\n", msg.City, msg.State)
	case msg := <-c2:
		fmt.Printf("Resultado encontrado: %v - %v (Fonte: BrasilAPI)\n", msg.City, msg.State)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}
}

func SearchInViaCep(cep string, ch chan Address) {
	request, err := http.Get("http://viacep.com.br/ws/" + cep + "/json")
	if err != nil {
		log.Printf("%v", err)
	}
	defer request.Body.Close()

	result, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("%v", err)
	}

	var data ViaCEP
	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Printf("%v", err)
	}

	// PARA TESTE
	// time.Sleep(time.Second * 2)

	// CASO NÃO ENCONTRE O CEP
	if data.Cep != "" {
		ch <- Address{
			Cep:          data.Cep,
			State:        data.State,
			City:         data.City,
			Neighborhood: data.Neighborhood,
			Street:       data.Street,
		}
	}
}

func SearchInBrasilApi(cep string, ch chan Address) {
	request, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		log.Printf("%v", err)
	}
	defer request.Body.Close()

	result, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("%v", err)
	}

	var data BrasilAPI
	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Printf("%v", err)
	}

	// PARA TESTE
	// time.Sleep(time.Second * 2)

	// CASO NÃO ENCONTRE O CEP
	if data.Cep != "" {
		ch <- Address{
			Cep:          data.Cep,
			State:        data.State,
			City:         data.City,
			Neighborhood: data.Neighborhood,
			Street:       data.Street,
		}
	}
}
