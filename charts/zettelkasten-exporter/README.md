# Instalação

Esta seção descreve como implantar o `zettelkasten-exporter` em um cluster [Minikube](https://minikube.sigs.k8s.io/docs/) local como parte do trabalho da disciplina de DevOps.

### 1. Iniciar o Minikube

Primeiro, inicie o seu cluster Minikube:

```bash
minikube start
```

### 2. Habilitar o Ingress

Habilite o addon de Ingress no Minikube para expor os serviços:

```bash
minikube addons enable ingress
```

### 3. Configurar o DNS local

Adicione uma entrada no seu arquivo `/etc/hosts` para apontar o DNS `k8s.local` para o IP do Minikube.

Primeiro, obtenha o IP do Minikube:

```bash
minikube ip
```

Em seguida, edite seu arquivo `/etc/hosts` (você precisará de permissões de administrador) e adicione a seguinte linha, substituindo `<MINIKUBE_IP>` pelo endereço obtido no passo anterior:

```
<MINIKUBE_IP> k8s.local
```

**Nota:** O `zettelkasten-exporter` não possui uma interface de usuário, então o Ingress principal será para o Grafana.

### 4. Instalar o Helm Chart

Instale o Helm chart a partir do registro OCI `ghcr.io`. Você pode usar o arquivo `values.sample.yaml` como ponto de partida para a sua configuração.

```bash
helm install zettelkasten-exporter oci://ghcr.io/luissimas/zettelkasten-exporter-chart/zettelkasten-exporter -f values.sample.yaml --namespace zettelkasten --create-namespace
```

Após a instalação, a aplicação estará rodando no seu cluster Minikube. Para acessar a aplicação, basta seguir as instruções exibidas no terminal após executar o comando `helm install`.

Acesse a lista de dashboards disponível em <http://k8s.local/dashboards> e selecione o dashboard "Zettelkasten". Após alguns segundos, deverá ser possível visualizar as métricas exportadas pela aplicação.
