# Imagem base para contêiner
FROM golang:latest

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos necessários para o diretório de trabalho
COPY . .

# Compila o código Go
RUN go build -o sample .

# Define o comando padrão para executar o seu aplicativo quando o contêiner for iniciado
CMD ["./sample"]