# このリポジトリは何か
モダンな開発環境の一例として、以下のシステム構成の環境を提供するサンプルプロジェクトです。  
本READMEに含まれる環境構築手順に従ってセットアップできます。

**システム構成図(全体)**  
![モダン開発環境-Lv2(具体)](https://github.com/user-attachments/assets/46b8aeeb-b753-41bc-a4af-cedcad2697a9)

システムの構成要素|説明
---|---
Gitホスティングサービス|GitHub
クラウドインフラ|GCP:<br>- Cloud Run(コンテナサービス)<br>- Artifact Registry(Dockerレジストリ)<br>- Spanner(DB)<br>
Webアプリケーション|- 機能: gRPCのAPI。サンプルデータの追加、一覧取得、更新、削除を行うことができる。<br>- 実装方法:<br>  - 言語: Golang<br>  - フレームワーク: 未使用<br>  - 主なライブラリ: ProtocolBuffers, GORM
CI/CDツール|CircleCI
連絡ツール|Slack
ローカルマシン|Mac

**システム構成(ローカルマシン)**
![モダン開発環境-ローカル環境のシステム構成、コンテナオーケストレーションの設定](https://github.com/user-attachments/assets/c5ea98bf-9b85-4238-98fb-f95a3bbf4937)


**システム構成(CircleCI, リモート環境)**
![モダン開発環境-Lv2(GCP, CircleCI)](https://github.com/user-attachments/assets/788b3230-7f19-4eb0-864f-bf07ddc2a7a7)

---
# 2. 環境構築手順
## 2.1. 各サービス、ツールの初期設定
各サービスについて、無料でゼロから準備する前提での手順となっています。  
すでに試すことのできるアカウント等をお持ちの場合、ツールをインストール済みの場合は、適宜手順をスキップしてください。  
また、各ツール、サービスとも経年変化により手順が多少変わる可能性があります。  

### 2.1.1. GitHub
#### GitHubのアカウントを作成
[GitHub](https://github.com/)へアクセス > Sign up
以降、画面の指示に従って作成してください
Free/Teamsは任意。Freeの場合かつ、Privateリポジトリにする場合は、後述のリポジトリの保護ルールが適用できません。
 
### 2.1.2. GCP
#### Googleアカウントを作成
GCPインフラを構築するGoogleアカウントを用意してください。

#### GCPのFree Trialを申請
[GCPのFree Trial(90日間$300で使い放題)](https://cloud.google.com/free?hl=ja)を申請してください。  
クレジットカード情報が必要になります。
 
#### gcloud CLIをインストール
[gcloud CLIをインストールする](https://cloud.google.com/sdk/docs/install?hl=ja)に従ってインストールしてください。

### 2.1.3. Go
<TODO>
 
### 2.1.4. CircleCI
#### CircleCIのアカウントを作成
①[CircleCI](https://circleci.com/)へアクセス > Sign up > Sign up  
②任意のメールアドレス、パスワードを設定してアカウントを作成  
③Start a new organization - Get Startedを押下  
![circleci-start-new-org](https://github.com/user-attachments/assets/7a5a03b6-e752-4282-91f0-87ae3c9eea24)

④任意の組織名を入力してLet's goを押下  
 
### 2.1.5. Slack
#### Slackのアカウントを作成
[Slack](https://slack.com)へアクセス してアカウントを作成してください。
 
#### Slackのワークスペース、チャンネルを作成
GitHubやCircleCIからの通知先となるワークスペース、チャンネルを作成してください。
 
### 2.1.6. Docker
[Docker Desktop](https://docs.docker.com/desktop/install/mac-install/)からDocker for Desktopをインストールしてください。


### 2.1.7. Terraform
#### Terraformをインストール
```
# tfenvをHomebrewでインストール
brew install tfenv
  
# インストール可能なterraformのバージョンを確認
tfenv list-remote
  
# terraformをバージョンを指定してインストール
tfenv install 1.9.1
  
# インストール済みのterraformのバージョンを一覧表示(*付きが現在使用しているバージョン)
tfenv list
  
# 使用するterraformのバージョンを切り替え
tfenv use 1.9.1 
```
実務上、terreformのバージョンアップ作業が発生する可能性を考慮して、tfenv経由でインストールしています。
 
 
## 2.2. 各サービス、ツールの詳細設定
**前提** 
[2.1. 各サービス、ツールの初期設定](#21-各サービス、ツールの初期設定) の作業を終えていることが前提です。
 
### 2.2.1. GitHub
#### 2.2.1.1. リポジトリの作成
GitHubへログイン > New repository > 以下の設定で Create repository
 
項目|値
---|---
Repository name|任意
Description|任意
Public/Private|Private
Add a README file|チェックしない
Add .gitignore|None
Choose a license|None
 
 
#### 2.2.1.2. サンプルプロジェクトをリポジトリに登録
[サンプルプロジェクト](https://github.com/fnami0316/modern-dev-env-app-sample)をコピーしたリモートリポジトリをGitHubに作成します。
 
##### 手順
ローカルマシンのターミナルで以下を実行してください。

SSHキーの生成
```
cd ~/.ssh
ssh-keygen -t ed25519 -C "<作成したGitHubアカウントのメールアドレス>"
```

SSHキー(公開鍵)をGitHubへ登録
GitHub(Web)へログイン > 右上のアバターアイコン > SSH and GPG keys > New SSH key で公開鍵をコピーして登録

SSHの設定ファイル編集
```
Host <任意のホスト名>
  HostName github.com
  User git
  IdentityFile "<SSH秘密鍵のパス>"
```

 
サンプルプロジェクトを自身のリモートリポジトリへコピー
```
# サンプルプロジェクトをcloneするディレクトリへ移動
cd <任意のディレクトリ>

# サンプルプロジェクトをclone
git clone https://github.com/fnami0316/modern-dev-env-app-sample.git

# サンプルプロジェクトのディレクトリへ移動
cd modern-dev-env-app-sample

# 使用するgitアカウントを設定
git config --local user.name "<任意の名前>"
git config --local user.email "<作成したGitHubアカウントのメールアドレス>"

# gitアカウントの設定確認
git config user.name
git config user.email

# ローカルブランチ(master)のリモートリポジトリを、サンプルプロジェクトのものから自身のGitHubアカウントのリモートリポジトリへ変更
git remote set-url origin git@<↑で設定した任意のホスト名>:<GitHubアカウント名>/<リポジトリ名>.git
 
# 現在のリモートリポジトリを確認
git remote -v
 
# 自身のGitHubアカウントのリモートリポジトリにサンプルプロジェクトをpush
git push origin
```
 
#### 2.2.1.3. terraformの設定を変更
以下変更して、commit & pushしてください。

①deploy/terraform/variables.tfを変更します。

```
variable "default_project_id" {
  type = string
  description = "各リソースのデフォルトのプロジェクトID"
  default = "変更必須"
}

variable "spanner_instance_dev" {
  description = "Cloud Spannerのインスタンスの設定値"
  type = object({
    name = string
    config = string
    display_name = string
    num_nodes = number
  })
  default = {
    # インスタンス名(ID)
    name = "変更必須"
    # インスタンス構成
    config = "regional-us-central1"
    # 表示名
    display_name = "変更必須"
    # ノード数
    num_nodes = 1
  }
}

```
- default_project_id の defaultを、自身のGCPプロジェクトID(xxxxxxxx-xxxxxx-xxxxxx-xxといったフォーマットのもの)に変更してください。プロジェクトIDは、GCPのWebコンソールから確認できます。
- spanner_instance_dev の name, display_nameを任意の値に変更してください。迷う場合、とりあえず"<GCPプロジェクト名の一部>-dev"といった形にしておくのが楽です。
- spanner_database_dev の instance を spanner_instance_dev の nameと同じものに、name を任意の値に変更してください。迷う場合、nameはとりあえず"<Spannerインスタンス名>-1"といった形にしておくと楽です。

②deploy/terraform/backend.tfを変更します。

```
terraform {
  backend "gcs" {
    bucket = "変更必須"
  }
}

# tfstateを保存するためのGCSバケット
resource "google_storage_bucket" "tfstate" {
  name     = "変更必須"
  location = var.default_region
  storage_class = "STANDARD"
  versioning {
    enabled = true
  }
}
```

terraform.backend.gcs.bucketと、google_storage_bucket.tfstate.name にあるtfstateファイルの保存先のGCSバケット名を任意の名前に変更します。
[GCSのバケット名はグローバルでユニークでなければならず、プロジェクトIDとかも入れない方がよい](https://pms-confluence2.banadev.com/confluence/pages/viewpage.action?pageId=1099910187#:~:text=%E4%BF%9D%E5%AD%98%E5%85%88%E3%81%AE%E3%83%90%E3%82%B1%E3%83%83%E3%83%88%E5%90%8D%E3%81%AF%E3%82%B0%E3%83%AD%E3%83%BC%E3%83%90%E3%83%AB%E3%81%A7%E3%83%A6%E3%83%8B%E3%83%BC%E3%82%AF%E3%81%A7%E3%81%AA%E3%81%91%E3%82%8C%E3%81%B0%E3%81%AA%E3%82%89%E3%81%9A%E3%80%81%E3%83%97%E3%83%AD%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88ID%E3%81%A8%E3%81%8B%E3%82%82%E5%85%A5%E3%82%8C%E3%81%AA%E3%81%84%E6%96%B9%E3%81%8C%E3%82%88%E3%81%84)とのことで、"tfstate-<適当なUUID>"といった形にするのが無難かと思います。
 


####  2.2.1.4. リポジトリの設定を変更(任意)
複数人で開発することを想定し、利便性向上、誤操作による復旧の手間を低減するという観点から、以下の項目を設定します。
- wiki利用可: デフォルト設定で可能
- masterブランチに対して以下の制限を設定(ブランチ保護ルール)
    - force pushを許可しない
    - 他のブランチをマージする前にPR必須にする
    - マージするブランチ側でのテスト必須にする
 

※Freeプランかつ、Privateリポジトリにする場合は、このリポジトリの保護ルールは適用できないため無視してください。
##### 手順
①GitHubへログイン > 対象のリポジトリのSettings > Code and automation - Rules - Rulesetsで New ruleset  > New branch ruleset
→以下の画面が表示されます。
 
![image2024-7-5_16-49-30-](https://github.com/user-attachments/assets/f386aa20-b93f-4e6c-89d5-e3f811fb0a05)

 
②各項目を以下の通り設定してCreate押下
項目|　|値|補足
---|---|---|---
RulesetName| |任意のブランチ保護ルール名|ブランチ保護ルールを管理する上での名前。"masterブランチ"などわかりやすいものを設定。
Enforcement status| |Atcive|このブランチ保護ルールを適用するか否か
Targets|Target branches|master|ブランチ保護ルールを適用するブランチ名のパターン。Include by patternから設定
Rules|Branch protections|- Require a pull request before merging: true<br>　- Required approvals: 1<br>　- Dismiss stale pull request approvals when new commits are pushed: true<br>　- Require approval of the most recent reviewable push: false<br>　- Require conversation resolution before merging: false<br>|- マージ前のPR必須にするか<br>　- 必須のapprove数<br>　- approve後にpushがあった場合、過去のapproveを却下する<br>　- 直近のレビュー可能なpushのapproveを必須にするか<br>　 - マージ前の会話の解決を必須にするか<br>
　|　|- Require status checks to pass: true<br>　- Require branches to be up to date before merging: true<br>　- ci/circleci: build-and-test|- マージしようとしているブランチのステータスチェック通過を必須にするか<br>　- 最新のコードでチェックしなければならないか(後述のCircleCIとGitHubの連携設定が終わった後でないと設定不可)<br>　 
　|　|Block force pushes: true|
 
##### 参考 
- [プランと請求日を表示する](https://docs.github.com/ja/billing/managing-your-github-billing-settings/viewing-your-subscriptions-and-billing-date)
- [ウィキについて](https://docs.github.com/github/building-a-strong-community/about-wikis)
- [プルリクエストを自動的にマージする](https://docs.github.com/ja/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request)
- [保護されたブランチについて](https://docs.github.com/ja/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches#restrict-who-can-push-to-matching-branches)

### 2.2.4. Terraformを使ってGCPのサービスをプロビジョニングする
#### 手順
ローカルマシンのターミナルで以下を実行してください。
 
```
# 1. terraform(Google Cloud SDK を使ったアプリケーション)のための認証
gcloud auth application-default login
  
# 2. ブラウザでGCPプロジェクトを操作するGoogleアカウントにログイン。ログイン後ターミナルを再び操作できるようになる。
# ブラウザのクッキーを全削除しておかないとうまく動作しないかも
  
# 3. terraformのルートモジュールへ移動
cd deploy/terraform
  
# 4. backend.tfのterraformブロックをコメントアウト
  
# 5. terraformで必要になるproviderをセットアップ
terraform init
  
# 6. dry run
terraform plan
  
# 7. 適用(tfstateの保存先はローカル)
terraform apply
  
# 8. backend.tfのterraformブロックをコメントアウトを元に戻す
  
# 9. dry run
terraform plan
  
# 10. 適用(tfstateの保存先がgcsに)
terraform apply
```

補足:  

run.googleaips.comを有効化するとArtifactRegistryのAPIも暗黙的に有効化されるため、Cloud Runのリソースでdepends_on = [google_project_service.gcp_services["run.googleapis.com"]] を施すことで必要なAPIの有効化 → Cloud Runサービスの作成という順序でのリソースの設定を期待しています。
しかし実際にはうまく機能しないようで、初回の実行ではArgifactRegistryのAPIが有効化されていないというエラーが出力されてしまうかもしれません。
これはgoogle providerの機能不備だと感じています(google providerに限らずよくあることです)。
無視して再度`terraform apply`を実行すると成功するはずです。

TODO: イメージが存在しないエラー
Cloud Runサービスのプロビジョニングでは、リビジョンで使用するイメージの設定が必須になります。
対象のイメージがGARに存在しない場合、以下のようなエラーが発生します。

```
Error: Error waiting to create Service: Error waiting for Creating Service: Error code 13, message: Revision 'api-00001-nps' is not ready and cannot serve traffic. Image 'us-central1-docker.pkg.dev/dazzling-pillar-435904-a5/api/sample_app:latest' not found.
│
│   with google_cloud_run_v2_service.api,
│   on main.tf line 18, in resource "google_cloud_run_v2_service" "api":
│   18: resource "google_cloud_run_v2_service" "api" {
```

そのため、ダミーのイメージを使ってデプロイ後、CircleCIのmasterブランチへのpushのパイプラインで正式なイメージをビルド&デプロイした後に、再度terraform applyで適用するというめんどくさい工程を取ります。


---
 
### 2.2.2. CircleCI
#### 2.2.2.1. CircleCIプロジェクトの作成とGitHub連携
①CircleCIへログイン > Projects > CreateProject  
![circleci_create_project](https://github.com/user-attachments/assets/6d14082a-2159-44b9-9224-a0c50c02bec8)

 
②Build, test, and deploy a software application を選択  
![image](https://github.com/user-attachments/assets/ded0f376-7465-4331-8b53-11ecdb1922d5)


③任意のプロジェクト名を入力して Next: set up a pipeline を押下  
![image](https://github.com/user-attachments/assets/8b1c0c30-b1d9-4792-b546-ddb1ce586996)


④任意のCircleCIプロジェクト名を 入力してNext: choosse a repo を押下  
![image](https://github.com/user-attachments/assets/89a7f97d-18e0-4656-9573-98a7821b4115)

⑤Connect to GitHubを選択  
![image](https://github.com/user-attachments/assets/3b3d9288-1378-4b4a-af7b-62c014068972)

⑥連携するGitHubのアカウントでログイン  
![image](https://github.com/user-attachments/assets/2d4a73f2-f8af-402e-8163-2cafe127e450)

(あらかじめGitHubからログアウトしておかないと、目標のアカウントとうまく連携できないことがあります)

⑦Only select repositories を選択し、連携するリポジトリを選択して Install & Authorize  
![image](https://github.com/user-attachments/assets/0093e93f-605c-489e-a09b-a3a43a5c27b1)

→ CircleCIのWebページに戻ってきます
![image](https://github.com/user-attachments/assets/80680992-9b80-4124-ba5a-89e017ab0f7c)

⑧連携するリポジトリを選択して、Next: set up your config を押下  
![image](https://github.com/user-attachments/assets/e55ccc88-9828-4e5f-a42d-806e6090cb02)

⑨Next: set up triggers を押下  
![image](https://github.com/user-attachments/assets/9a1fffa1-6bf7-4bd4-8908-cc9863178aab)

⑩Next: review and finish setup を押下  
![image](https://github.com/user-attachments/assets/4b06ceac-adf4-4f35-96a1-18c2134fa309)

⑪Finish setup を押下  
![image](https://github.com/user-attachments/assets/e376adcc-2736-497d-a647-5561e2dacc04)

---

#### 2.2.2.2. CircleCIプロジェクトとSlackを連携
##### 通知先となるSlackチャンネルのワークスペースにSlackアプリを作成
**①アプリの作成** 
Slack API Webサイトへアクセス > Your Apps > ログイン > Create New App > From scratch > App Name に任意のアプリ名を入力して対象のワークスペースを選択　> Create App
 
**②アプリの権限設定** 
Settings - OAuth & Permissions > Scopes - Bot Token Scopes > Add an OAuth Scope 押下後、以下の選択を繰り返す
- chat:write
- chat:write.public
- files:write

![image](https://github.com/user-attachments/assets/109ff27c-0bbe-4e48-a9ed-43ca8579aad1)
![image](https://github.com/user-attachments/assets/5e48fe4d-44ad-4108-9855-6283049b33e2)


 **③アプリのインストール**
Features-OAuth & Permissions > Install to <ワークスペース名> > 許可する  
→ OAuth Tokens - Bot User OAuth Token に トークンが表示される (次の工程で使うのでコピーしておく)
 
 
#### 2.2.2.3. CircleCIプロジェクトに環境変数を設定
CircleCIへログイン > 事前に作成した組織を選択 > Projects > 事前に作成したプロジェクトを選択 > Project Settings > Environment Variables より、それぞれ登録してください。

##### Slack連携関連
slack Orbが使用します。

**Slackアプリのアクセストークン**
 
項目|値
---|---
Environment Variable Name|SLACK_ACCESS_TOKEN
Value|Slack API のWebサイトで作成したアクセストークン

**Slackチャンネル**
 
項目|値
---|---
Environment Variable Name|SLACK_DEFAULT_CHANNEL
Value|Slackのデスクトップアプリ上で、通知先のSlackチャンネルを右クリック > コピー > リンクをコピーで得られるURLの末尾の文字列(多分11文字)

##### GCP操作関連
gcp-cli orbが使用します。
GARへのDockerイメージのPUSHやCloud Runへのデプロイ処理で必要になります。

**GCPサービスアカウントのキー**
 
項目|値
---|---
Environment Variable Name|GCLOUD_SERVICE_KEY
Value|GCPコンソールで作成したサービスアカウントのキー

**GCPプロジェクトID**
 
項目|値
---|---
Environment Variable Name|GOOGLE_PROJECT_ID
Value||.tfにセットしたXXXX(TODO)の値

**GCPリージョン**
 
項目|値
---|---
Environment Variable Name|GOOGLE_COMPUTE_REGION
Value||.tfにセットしたXXXX(TODO)の値

**GCPゾーン**
 
項目|値
---|---
Environment Variable Name|GOOGLE_COMPUTE_ZONE
Value|.tfにセットしたXXXX(TODO)の値


#### 2.2.3.3. 動作確認
サンプルプロジェクトをコピーした自身のGitHubリポジトリへ空コミットでpushします。
```
git commit --allow-empty -m "空コミット"
git push origin HEAD:master
```

設定が正しく行えていれば、Slackチャンネルに以下のようなWorkflowの開始通知が届くはずです。
![スクリーンショット 2024-09-18 15 05 08](https://github.com/user-attachments/assets/8fe8f68f-ba47-4780-afd2-36aef36fea0c)

また、すべての工程が成功すると以下のような成功通知が届きます。
masterブランチへのpushの場合、Cloud Runへのデプロイまで行われます。
CircleCIでデプロイした場合、デプロイしたリビジョンにcircleci-deployタグが付与されます。
タグ付きURLを利用して、以下のコマンドで動作確認できます。

```
## サンプルデータを追加
grpcurl  -d '{"name": "sample1"}' circleci-deploy---xxx.a.run.app:443 api.SampleService.CreateSample
## サンプルデータを更新
grpcurl -d '{"id": "<サンプルデータ追加時のレスポンスに含まれるID>", "name": "updated-sample1"}' circleci-deploy-xxx.a.run.app:443  api.SampleService.UpdateSample
## サンプルデータの一覧を取得
grpcurl -d '{"ids":["<サンプルデータ追加時のレスポンスに含まれるID>"]}' circleci-deploy-xxx.a.run.app:443 api.SampleService.ListSamples
## サンプルデータを削除
grpcurl -d '{"id":["<サンプルデータ追加時のレスポンスに含まれるID>"]}' circleci-deploy-xxx.a.run.app:443 api.SampleService.DeleteSample
```

xxxは、Cloud Runで発行されたURLからhttps://を削除したものになります。




### 2.2.5. ローカル環境を起動する
#### 手順
①Docker for Desktopを起動します。  
②ローカルマシンのターミナルで以下を実行してください。

```
# 1. プロジェクトルートへ移動
cd <プロジェクトルート>

# 2. ローカル環境を起動
make docker-compose-up-d

# 3. 動作確認(Goアプリケーションソースのテスト)
# 全て?もしくはokならば、正常に動作しています。
make docker-go-test-serial

# 4. 動作確認(APIへリクエスト)
# ランダムなidでレスポンスが返ってきたら正常に動作しています。
grpcurl -plaintext -d '{"name": "sample1"}' localhost:8080 api.SampleService.CreateSample
```





