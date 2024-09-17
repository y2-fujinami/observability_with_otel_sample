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
 
### 2.2.1. GitHubリポジトリの作成と設定
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
 
####  2.2.1.3. リポジトリの設定を変更
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
 
 
### 2.2.2. GitHub連携したCircleCIプロジェクトを作成
#### 手順
①CircleCIへログイン > Projects > CreateProject  
![circleci_create_project](https://github.com/user-attachments/assets/6d14082a-2159-44b9-9224-a0c50c02bec8)

 
②Build, test, and deploy a software application > 任意のプロジェクト名を入力して Next: set up a pipeline
![circleci_what_would_you_like](https://github.com/user-attachments/assets/1c325039-9971-40be-a32a-1dd1ff5d6b9f)

 
③GitHubを選択  
![circleci_select_repo](https://github.com/user-attachments/assets/99bd6f0f-df09-4c79-b1ed-12097bfc4982)


④連携するGitHubリポジトリのアカウントを選択 > Only select repositories にチェックを入れて、連携するリポジトリを選択して Install & Authorize  
![circleci_install_authorize](https://github.com/user-attachments/assets/d1e4c380-d102-4416-a7d1-746c4671d50e)

⑤Project Name へ任意のCircleCIプロジェクト名を入力

⑥Use Existing Configを押下  
![use_existing_config](https://github.com/user-attachments/assets/4cb83562-84b4-4a17-af29-cf0efb6e9d15)

⑦Start Building を押下  
![circleci-start-building](https://github.com/user-attachments/assets/05ba9085-2696-4d1a-87b5-1c40fc3ecb65)

この段階では、CircleCIとSlackを連携する設定をしていないため、Slack通知のjobでエラーになってしまいます。  
![circleci-error](https://github.com/user-attachments/assets/7a8b065a-95ee-4dd7-bddc-062275b885ea)

[2.2.3. CircleCIとSlackを連携](#223-CircleCIとSlackを連携)で設定します。
 
### 2.2.3. CircleCIとSlackを連携
#### 2.2.3.1. 通知先となるSlackチャンネルのワークスペースにSlackアプリを作成
##### 手順
**①アプリの作成** 
Slack API WebサイトのYour Apps > Create an App > From scratch > App Name に任意のアプリ名を入力して対象のワークスペースを選択　> Create App
 
**②アプリの権限設定** 
Settings-Basic Information > Add features and functionality > Permissions > Scopes - Bot Token Scopes > Add an OAuth Scope 押下、以下の選択を繰り返す
- chat:write
- chat:write.public
- files:write
 
**③アプリのインストール**
Features-OAuth & Permissions > Install to Workspace > 許可する > Slackアプリのインストール申請が承認された後、Features-OAuth & Permissions > Install to Workspace > 許可する
→ Bot User OAuth Token に トークンが表示される。次の工程で使うのでメモ。
 
 
#### 2.2.3.2. CircleCIの環境変数にSlackアプリのアクセストークンと通知先チャンネルを登録
##### 手順
CircleCIへログイン > [CircleCIのアカウントを作成](#CircleCIのアカウントを作成) で作成した組織を選択 > Projects > [2.2.2. GitHub連携したCircleCIプロジェクトを作成](#222-GitHub連携したCircleCIプロジェクトを作成)で作成したプロジェクトを選択 > Project Settings > Environment Variables より、以下2つの環境変数をそれぞれ登録
 
**Slackアプリのアクセストークン**
 
項目|値
---|---
Environment Variable Name|SLACK_ACCESS_TOKEN
Value|[2.2.3.1. 通知先となるSlackチャンネルのワークスペースにSlackアプリを作成](#2231-通知先となるSlackチャンネルのワークスペースにSlackアプリを作成)で作成したアクセストークン
 
**Slackチャンネル**
 
項目|値
---|---
Environment Variable Name|SLACK_DEFAULT_CHANNEL
Value|Slackのデスクトップアプリ上で、通知先のSlackチャンネルを右クリック > コピー > リンクをコピーで得られるURLの末尾の文字列(多分11文字)
 
#### 2.2.3.3. 動作確認
[2.2.2. GitHub連携したCircleCIプロジェクトを作成](#222-GitHub連携したCircleCIプロジェクトを作成)で失敗していたbuild-and-testワークフローを再実行して、Slackチャンネルに通知が届くことを確認する。  
![circleci-notification](https://github.com/user-attachments/assets/66240bf5-1478-4db1-9e2f-d6c8f206f9e7)

 
### 2.2.4. Terraformを使ってGCPのサービスをプロビジョニングする
#### 手順
ローカルマシンのターミナルで以下を実行してください。
 
```
# 1. terraform(Google Cloud SDK を使ったアプリケーション)のための認証
gcloud auth application-default login
  
# 2. ブラウザでGCPプロジェクトを操作するGoogleアカウントにログイン。ログイン後ターミナルに戻ってくる。
# ブラウザのクッキーを全削除しておかないとうまく動作しないかも
  
# 3. terraformのルートモジュールへ移動
cd deploy/terraform
  
# 4. terraform.backendをコメントアウト
  
# 5. terraformで必要になるproviderをセットアップ
terraform init
  
# 6. dry run
terraform plan
  
# 7. 適用(tfstateの保存先はローカル)
terraform apply
  
# 8. terraform.backendコメントアウトを元に戻す
  
# 9. dry run
terraform plan
  
# 10. 適用(tfstateの保存先がgcsに)
terraform apply
```

### 2.2.5. ローカル環境を起動する
#### 手順
ローカルマシンのターミナルで以下を実行してください。

```
# 1. プロジェクトルートへ移動
cd <プロジェクトルート>

# 2. ローカル環境を起動
make docker-compose-up-d

# 3. 動作確認(Goアプリケーションソースのテスト)
make docker-go-test

# 4. 動作確認(APIへリクエスト)
grpcurl -d '{"name": "sample1"}' localhost:8080 api.SampleService.CreateSample
```



