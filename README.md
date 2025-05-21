# temperatures
Projeto de estudos, consulta CEP e retorna temperatura - Go Expert

## Descrição
São duas APIs onde a API-A espera um POST, que por sua vez faz uma request na API-B que faz uma consulta que espera o cep e retorna a temperatura da cidade a qual o cep correponde, o retorno da API traz a temperatura em Celsius, Fahrenheit e Kelvin.

Para isso, utilizei duas API:

- [viacep](https://viacep.com.br/) que permite buscar um cep, a API é gratuita;
- [weatherapi](https://www.weatherapi.com/) que traz dados de uma cidade, essa API é necessário um cadastro;

## Funcionamento
A primeira API-A espera um POST e faz outra request na API-B GET com o parâmetro "zipcode", ao receber a request é feita uma validação para filtrar erros de digitação como uma letra no lugar de algum número no cep, em seguida é feita a consulta no **viacep**, e com a cidade retornada é feita uma consulta no **weatherapi** que retorna vários dados dentre eles a temperatura em celcius e fahrenheit, ficando apenas necessário o calculo da temperatura em kelvin, após isso o sistema retorna um JSON com as três temperaturas.

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
curl -X POST http://localhost:8080 -H "Content-Type: application/json" -d '{"zipcode": "11390300"}'
```

## OTEL com Zipikn

Um dos requisitos era termos um colletor OTEL que envie os dados para o Zipkin, as configurações podem ser consultadad no arquivo "otel-collector-config.yml" que fica dentro da pasta docker, e após fazer uma request POST como mostrado acima é feito o trakeamento das requests como é mostrado no print abaixo:

![image](https://github.com/user-attachments/assets/40fb4a2f-1e0c-4544-85a8-ea5de980d276)
