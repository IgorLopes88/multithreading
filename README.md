# Go Expert

Desafio **Multithreading** do curso **Pós Go Expert**.

### Execução
Para executar a aplicação e realizar a consulta do **CEP**, na pasta raiz, execute o comando abaixo, seguido do cep pretendido.
```
go run main.go 13330-250
```
Se digitar um **CEP VÁLIDO**, a resposta deverá ser a cidade e o estado, seguido de qual API conseguiu fornecer a resposta mais rápida.
Como mostrada abaixo:

```
Resultado encontrado: Indaiatuba - SP (Fonte: ViaCep)
```

ou

```
Resultado encontrado: Indaiatuba - SP (Fonte: BrasilAPI)
```


Caso o tempo de resposta seja superior a 1 segundo. A aplicação retornará com erro:

```
Timeout
```

Pronto!