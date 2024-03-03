## GoExpert

### Para a execução da aplicação e seus tests, basta executar os seguintes comandos
#### Execução local
Para executar a aplicação localmente, primeiramente você deve definir a API KEY do Weather API.
PS: Se você não tem uma conta, pode definir uma aqui: https://www.weatherapi.com/

```sh
export WEATHER_API_TOKEN=<API_KEY>
```

Após isso para inicializar a aplicação é muito simples, basta você executar o seguinte comando.

```sh
docker-compose up app --build 
```

Para exeutar os testes, basta executar os comandos a seguir

```sh
docker-compose up app-tests --build 
```

Isso estará funcionando na porta :8080 do seu localhost.

PS: Com relação ao endpoint de serviço, ainda não disponibilizei porque não tenho mais freetier na AWS e GCP, deixei uma dúvida no canal administrativo do discord, pergunta se posso subir em outros serviços.

### Documentação de Endpoint
URL /weather/:cep
parametro: CEP da sua localidade

