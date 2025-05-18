# temperatures
Projeto de estudos, consulta CEP e retorna temperatura - Go Expert

## Descrição
Uma api que espera o cep e retorna a temperatura da cidade a qual o cep correponde, o retorno da API traz a temperatura em Celsius, Fahrenheit e Kelvin.

Para isso, utilizei duas API:

- [viacep](https://viacep.com.br/) que permite buscar um cep, a API é gratuita;
- [weatherapi](https://www.weatherapi.com/) que traz dados de uma cidade, essa API é necessário um cadastro;

## Funcionamento

A API espera um GET com o parâmetro "zipcode", ao receber a request é feita uma validação para filtrar erros de digitação como uma letra no lugar de algum número no cep, em seguida é feita a consulta no **viacep**, e com a cidade retornada é feita uma consulta no **weatherapi** que retorna vários dados dentre eles a temperatura em celcius e fahrenheit, ficando apenas necessário o calculo da temperatura em kelvin, após isso o sistema retorna um JSON com as três temperaturas.

## Testes

Para fazer o teste, primeiro clone o projeto e em seguida é necessário criar uma conta no **weatherapi** para obter uma chave que será enviada nas requisições da api, em seguida renomeie o arquivo .env-eample, e coloquei sua chave nesse arquivo na variavel **WEATHER_API_KEY** para example, builde o dockerfile e execute o container:

```bash
mv .env-example .env
echo "AQUI SUA CHAVE GERADA LA NO WATHERAPI" >> .env
docker build --no-cache -t gotemp .
docker run --env-file .env gotemp
```

Dessa forma o container deve ficar rodando já esperando requisições. Para executar uma requisição faça um curl simples no terminal:

```bash
curl -X GET http://172.17.0.2:8080?zipcode=11703100
```