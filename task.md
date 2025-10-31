# GoでAPIを作成したい。
## 構成
まずはapiとresourcesに分けて欲しい。
resourcesの中も細分化してディレクトリ分けて欲しい。

## docker composeで起動するようにしたい。
docker composeはMakefileで操作できるようにする。
Makefileだけはリポジトリ直下に配置したい。

## APIについて
mysqlとredisとNewRelicが入っている状態。
mysqlで取得したデータをredisがキャッシュする構成。
また、redisはリクエストとその結果もキャッシュするようにして、高速に返却する。

NewRelicはAPMを見れればいいから、mysqlとredisにも注入して欲しい。

goのバージョンは1.25.3を使って欲しい。

また、open apiを使って定義して欲しくて、それはresourcesに配置してほしい。
open apiの定義書からコードを自動で生成し、それをGo側で使って欲しい。

Goのフレームワークとしてはechoを使って欲しい。

また、クリーンアーキテクチャで書いて

