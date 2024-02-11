# rest-api-simulator-service

Este projeto possibilita gerar um endpoint local para teste de conexão e tratamento de respostas de APIs RESTful, permitindo que configure diversas rotas e métodos direferentes, permitindo indicar o statuscode da resposta e o conteúdo.

Para utilizar é bastante simples, tudo que precisa é configurar o arquivo config.yaml com as rotas que deseja criar para o teste:

```yaml
TemplateFormatVersion: 2024-02-10
Description: config sample of lowcode-lambda with go

Server:
  Port: 8080

Resources:
  WhatsAppAuth:
    ResourceType: Endpoint
    Description: Simulação da autenticação do WhatsApp Business API
    Properties:
      Path: /v19.0/oauth/access_token
      Method: POST
      Response:
        StatusCode: 200
        Content: responses/auth.json

  WhatsAppMessage:
    ResourceType: Endpoint
    Description: Simulação do envio de mensagens do WhatsApp Business API
    Properties:
      Path: /v19.0/:phone_number/messages
      Method: POST
      Response:
        StatusCode: 200
        Content: responses/message.json
```

Na sessão _Server_ você deve indicar a porta que será usada para responder as requisições realizadas pelas aplicações cliente.

Já na sessão _Resources_ você irá indicar os recursos/rotas disponíveis, contendo:
- Tipo de resource: Endpoint
- Descrição (opcional)
- Propriedades:
  Path da API (indicar variáveis path com :)
  Método que será utilizado (POST, GET, DELETE, PATCH, PUT, OPTION ...)
- Response:
  Código do status que deseja retornar na requisição
  Caminho do arquivo JSON com a resposta da requisição recebida

É possível utilizar este projeto para criar cenários de erro e falhas e validar o comportamento da aplicação localmente.

Caso tenha alguma sugestão que possa ajudar a melhorar o projeto, favor abrir uma issue para que possamos avaliar.

