package mongodbproject

import (
	"bytes"
	"context"
	"fmt"
	"github.com/helmutkemper/chaos/factory"
	"github.com/helmutkemper/chaos/internal/manager"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime/debug"
	"testing"
	"time"
)

// Este é um exemplo de como instalar um container baseado em imagem e foi escrito seguindo a regra do programador cansado.

// No primeiro exemplo, você encontrará informações básicas para suas necessidades, porém, no exemplo abaixo, será adicionada mais informações para ser usada em um projeto complexo

// TestSimpleLinearBasic cria uma instalação simples do banco de dados MongoDB e a deixa funcionar por três minutos.
//
// Nesse exemplo serão mostradas as configurações básicas de como criar um container baseado em imagem e como expor uma porta ao mundo
func TestSimpleLinearBasic(t *testing.T) {

	// Cria toda a infraestrutura necessária para projeto funcionar de forma adequada.
	primordial := factory.NewPrimordial().
		// [opcional] Cria uma rede dentro do docker, isolando o teste.
		//NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		// [opcional] Permite ao controlador de lixo remover qualquer imagem, rede, volume ou container criado para o teste
		// [opcional] "mongo:latest" remove a imagem ao final do teste, limpando espaço em disco
		//            como regra, todos os elementos criados pelo teste contém a palavra `delete` como um identificador de
		//            algo criado para o teste, porém, você pode passar nomes de elementos docker criados para o teste que
		// serão removidos ao final do teste
		Test(t, "mongo:latest")

	// Fábrica de container baseado em uma imagem existente
	factory.NewContainerFromImage(
		"mongo:latest",
	).
		// [opcional] Determina uma ou mais portas a serem expostas na rede
		Ports("tcp", 27017, 27017).
		// [opcional] Segue o manual do MongoDB e libera conexão de qualquer endereço
		EnvironmentVar([]string{"--bind_ip_all"}).

		// [opcional] permite gravar os dados do MongoDB em pasta local
		// Volumes("/data/db", "./data/db").

		// [opcional] permite reescrever o arquivo de configuração baseado em arquivo contido em pasta local
		// Volumes("/data/configdb/mongod.conf", "./data/configdb/mongod.conf").

		// [opcional] Espera pelo aparecimento de um texto na saída padrão do container, antes de prosseguir com o código
		WaitForFlagTimeout("Waiting for connections", 30*time.Second).
		// Determina o nome do container e a quantidade de containers a criada
		Create("mongo", 1).
		// Inicializa o container
		Start()

	// Quando a saída padrão do container imprimir o texto "Waiting for connections" o código chegará nesse ponto
	// Nesse momento, no diretório do projeto haverá os seguintes arquivos:
	//   report.mongo:latest.md: Relatório de segurança baseado no projeto https://github.com/google/osv-scanner
	//   stats.delete_mongo.0.csv: Relatório de consumo de memória e desempenho do container, baseado em capturas de dados pontuais

	// Caso queira controlar o tempo total do teste, crie uma go routine e deixe o teste ocorrer em paralelo
	go func(t *testing.T, primordial *manager.Primordial) {
		var err error
		var mongoClient *mongo.Client
		var start = time.Now()

		// Cria o cliente MongoDB
		mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://0.0.0.0:27017"))
		if err != nil {
			panic(string(debug.Stack()))
		}

		// Conecta ao MongoDB
		err = mongoClient.Connect(context.Background())
		if err != nil {
			panic(string(debug.Stack()))
		}

		// Testa a conexão
		err = mongoClient.Ping(context.Background(), readpref.Primary())
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			panic(string(debug.Stack()))
		}

		// Cria uma estrutura de dados
		type Trainer struct {
			Name string
			Age  int
			City string
		}

		// Insere os dados
		var collection *mongo.Collection
		var totalOfInserts int64 = 1000
		for i := int64(0); i != totalOfInserts; i += 1 {
			collection = mongoClient.Database("test").Collection("trainers")
			ash := Trainer{"Ash", 10, "Pallet Town"}
			_, err = collection.InsertOne(context.Background(), ash)
			if err != nil {
				panic(err)
			}
		}

		// Testa a integridade
		var total int64
		if total, err = collection.CountDocuments(context.Background(), bson.M{"name": "Ash", "age": 10, "city": "Pallet Town"}); err != nil {
			panic(err)
		}

		if total != totalOfInserts {
			t.Logf("total of inserts must be %v found %v", totalOfInserts, total)
			t.Fail()
		}

		fmt.Printf("fim\n")
		duration := time.Since(start)
		fmt.Printf("Duration: %v\n\n", duration)

		// [opcional] caso queira encerrar o teste antes do tempo determinado
		primordial.Done()
	}(t, primordial)

	// Determina o tempo do teste
	if !primordial.Monitor(3 * time.Minute) {
		t.Fail()
	}
}

// Este é um teste com criação de replica set para MongoDB.
// Caso tenha pulado a explicação anterior, ela contém o conhecimento básico de uso do sistema. Aqui são adicionadas mais informações
//
// Neste exemplo serão mostradas as configurações de comandos iniciais do container, acesso ao terminal, tratamento de resposta do terminal e conceitos de rede docker
func TestSimpleLinearComplex(t *testing.T) {

	primordial := factory.NewPrimordial().
		// Segundo manual do MongoDB, a replica set só ira funcionar se o hostname for definido, não podendo usar endereço IP
		// Para que o hostname funcione de forma correta, dentro do docker, deve ser criada uma rede, por isto, não comente
		// a criação de rede
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		// Caso queira continuar usando a imagem "mongo:latest", apenas não coloque o nome dela aqui
		Test(t)

	mongoDocker := factory.NewContainerFromImage(
		"mongo:latest",
	).
		// A função Create() manda criar 3 containers, por isto, o primeiro container terá a porta 27017 direcionada para
		// a porta 27016 na rede e assim por diante.
		// Caso seja passada apenas uma porta, apenas o primeiro container terá sua porta exposta na rede
		// Caso o container tenha várias portas, repita uma linha por porta
		Ports("tcp", 27017, 27016, 27017, 27018).

		// Imagine que você necessita usar a porta 27018, a porta usada para replicação secundária. Repita o comando com a próxima porta.
		// Ports("tcp", 27018, 27019, 27020, 27021).

		// Quando trabalhar com o MongoDB, tome cuidado na hora de especificar o IP a ser liberado, pois, bindIp: 0.0.0.0
		// irá liberar a conexão pelo IP 0.0.0.0, já o flag "--bind_ip_all" irá liberar para todos os IPs, o que é bem parecido, porém, NÃO É A MESMA COISA.
		// bindIp: 0.0.0.0 apenas aceita conexão especificada como sendo IP 0.0.0.0, --bind_ip_all aceita IP especificado como sendo 127.0.0.1, por exemplo.
		EnvironmentVar([]string{"--bind_ip_all"}).

		// Imagine que apenas o container 0 vai receber acesso externo
		// A lógica é simples: se apenas um valor é passado, ele serve para todos os containers, se mais de um valor é passado a chave 0 vai para container 0, a chave 1 vai para o container 1 e assim por diante
		// EnvironmentVar([]string{"--bind_ip_all"}, []string{}, []string{}).

		// Quando o container inicia, o MongoDB necessita receber o flag "--replSet NAME_REPLICA_SET", como no exemplo abaixo
		// $ docker run -p 27017:27017 --name mongo --net mongo_network mongo mongod --replSet rs0
		// No caso do MongoDB, alguns tutoriais omitem o nome do shell que irá receber o flag "--replSet rs0", dessa forma:
		// $ docker run -p 27017:27017 --name mongo --net mongo_network mongo --replSet rs0
		// Porém, no nosso caso, ele deve ser passado, como no comando abaixo
		Cmd([]string{"mongod", "--replSet", "rs0"}).

		// Caso necessite esperar por um flag indicador de sucesso, adicione um texto e o sistema ficará parado esperando pelo mesmo, porém, tome cuidado, e adicione o texto pensando em texto caso sensitivo.
		// Esta função usa a contains(str, text) para procurar texto, por isto, tome cuidado com texto muito curto
		WaitForFlagTimeout("Waiting for connections", 30*time.Second).

		// [opcional] Procura por textos indicadores de falha na saída padrão do container e salva a saída padrão do container na pasta indicada, para análise posterior. (cuidado: o texto é caso sensitivo)
		// Caso queira fazer um teste rápido, use a palavra "Waiting" e veja a pasta "bug" quando o container começar a rodar
		// Cuidado, ele procura por contains(str, text), por causa disso, o flag "fail" poderá encontrá as palavras "maxFailedInitialSyncAttempts" ou "failed". Por isto, eu coloco os dois pontos (:)
		FailFlag("./bug", "Address already in use", "panic:", "bug:").

		//Volumes("/etc/mongod.conf", "/Users/kemper/go/projetos/chaos/examples/mongoDbProject/conf/mongod_0.conf", "/Users/kemper/go/projetos/chaos/examples/mongoDbProject/conf/mongod_1.conf", "/Users/kemper/go/projetos/chaos/examples/mongoDbProject/conf/mongod_2.conf").

		// Determina a criação de 3 containers nos endereços 10.0.0.2:27016, 10.0.0.3:27017, 10.0.0.4:27018, host names delete_mongo_0, delete_mongo_1 e delete_mongo_2
		// Caso necessite mudar o host name, use a função HostName() e especifique um nome para cada container. Lembre-se, hostname requer uma rede anexada ao container
		Create("mongo", 3).

		// Embora seja obvio, é bom lembrar que as funções Create() e Start() devem ser as duas últimas funções chamadas
		Start()

	// Nesse ponto do código, os bancos estão prontos para uso, porém, a criação de replicas requer comandos do terminal
	// Para fazer isto, especifique a chave do container, 0 para o primeiro container criado o interpretador de comandos a ser usado e o flag indicador de que esses comandos serão enviados via texto, "-c", ou seja:
	// `/bin/bash -c "echo Hello World!"`

	// Escrever comandos de terminal para que o mongodb se transforme em replicaset
	var stdOutput []byte
	var err error

	// Para transformar o MongoDB do container de chave 2, delete_mongo_2, em secundário de replica set, é necessário acessar o container pelo terminal, acionar o terminal do MongoDB e passar o comando "rs.secondaryOk_()"
	// Explicação:
	//   * "/bin/bash": é o interpretador de comandos do linux
	//   * "-c": o comando vai chegar via string de texto, exemplo: `/bin/bash -c "echo Hello World!"`
	//   * "mongosh": é o interpretador de comandos do MongoDB
	//   * "127.0.0.1:27017": é o endereço de conexão na rede. Nesse ponto, perceba, o comando está acessando diretamente o container, e dentro do container, a porta é 27017 e o endereço é localhost. Não confunda o acesso interno, diretamente no container com acesso externo.
	//   * "--eval \"rs.secondaryOk_()\"": eval permite executar um comando javascript via texto, e como é texto dentro de texto, as aspas estão escapadas.

	// Quando isto acontecer, o comando vai devolver um texto contendo um indicador de erro ou sucesso em caso de falta de indicador. Os indicadores são:
	//   * DeprecationWarning: No MongoDB 6.0.6 pode ser ignorado
	//   * MongoNetworkError: Falha de conexão com o banco de dados
	//   * TypeError: Erro de sintaxe

	_, _, stdOutput, _, err = mongoDocker.Command(2, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.secondaryOk()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 2: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// Repete o mesmo processo para o container de chave 1, delete_mongo_1
	// Nota: Os containers secundários devem receber o comando "rs.secondaryOk()" antes do container principal receber o comando "rs.initiate()"
	_, _, stdOutput, _, err = mongoDocker.Command(1, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.secondaryOk()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 1: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// Inicializa o container de chave 0, delete_mongo_0, como sendo a instância MongoDB arbitro de replica set
	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.initiate()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// A falta de indicador de erro, por sí só é um indicador de sucesso, mas, o comando devolverá um json no seguinte formato:
	// {
	//   info2: 'no configuration specified. Using a default configuration for the set',
	//   me: 'delete_mongo_0:27017',
	//   ok: 1
	// }
	//
	// Caso você necessite processar o json, você poderá expressão regular, https://regex101.com/library/sjOfeq?orderBy=MOST_POINTS&page=3&search=json

	// Adiciona o MongoDB contido no container delete_mongo_1 como sendo membro do replica set
	// Notas:
	//   * Como o comando é passado via texto dentro de texto, cuidado com as aspas escapadas;
	//   * Dentro da rede docker, todos os MongoDB estão da porta 27017, as portas 27016, 27017 e 27018 são as portas expostas ao mundo, não na rede docker;
	//   * O host name "delete_mongo_x" só funciona dentro da rede docker
	//   * O MongoDB não aceita configuração de replica set por IP, apenas por host name
	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.add(\\\"delete_mongo_1:27017\\\")\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// Por pura curiosidade, em caso de sucesso, o MongoDB devolve o json:
	// {
	//   ok: 1,
	//   '$clusterTime': {
	//     clusterTime: Timestamp({ t: 1685399709, i: 1 }),
	//     signature: {
	//       hash: Binary(Buffer.from("0000000000000000000000000000000000000000", "hex"), 0),
	//       keyId: Long("0")
	//     }
	//   },
	//   operationTime: Timestamp({ t: 1685399709, i: 1 })
	// }

	// Adicione a próxima instância MongoDB ao replica set
	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh 127.0.0.1:27017 --eval \"rs.add(\\\"delete_mongo_2:27017\\\")\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if bytes.Contains(stdOutput, []byte("MongoNetworkError:")) || bytes.Contains(stdOutput, []byte("TypeError:")) {
		t.Logf("container 0: rs.secondaryOk().error: %s", stdOutput)
		t.FailNow()
	}

	// Caso você queira fazer uma última verificação, passe o comando "rs.status()" para qualquer instância mongo, ela deve devolver um json contando "set: 'rs0'" e "name: 'delete_mongo_x:27017'" para cada instância MongoDB
	_, _, stdOutput, _, err = mongoDocker.Command(0, "/bin/bash", "-c", "mongosh --eval \"rs.status()\"")
	if err != nil {
		t.Logf("mongoDocker.Command().error: %v", err.Error())
		t.FailNow()
	}

	if !bytes.Contains(stdOutput, []byte("'rs0'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_0:27017'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_1:27017'")) || !bytes.Contains(stdOutput, []byte("'delete_mongo_2:27017'")) {
		t.Logf("replica set, setup failed")
		t.FailNow()
	}

	// Nesse ponto do projeto, a replica set de MongoDB foi configurada com dados efêmeros e está em uma rede docker, comas portas 27016, 27017 e 27018 expostas ao mundo, mas, a replica set, por regra do MongoDB, só aceita conexão via host name, e host name só funciona na rede docker, por isto o teste deve ser feito em container

	// Cria um container a partir de uma pasta local
	factory.NewContainerFromFolder(
		"folder:latest",
		"./mongodbClient",
	).

		// Monta o dockerfile de forma automática caso o arquivo "main.go" esteja na raiz do projeto e o arquivo "go.mod" exista, mesmo que em branco.
		// Você pode especificar o caminho do Dockerfile, caso ele não esteja na raiz do projeto com o comando DockerfilePath("./path/inside/container/Dockerfile")
		MakeDockerfile().
		WaitForFlagTimeout("container is running", 10*time.Second).
		FailFlag("./bug", "panic:").
		Create("mongodbClient", 1).
		Start()

	// Deixa o projeto rodando por 5 minutos
	if !primordial.Monitor(5 * time.Minute) {
		t.Fail()
	}
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
