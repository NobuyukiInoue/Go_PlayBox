# Go_PlayBox

MIDI関連のWIN32APIをGolangから呼び出して楽譜データを演奏します。 （MS-Windowsでしか動作しません。）

qiita.comにて記事を掲載しています。

リンク）  
[Golangで演奏する電子オルゴール](https://qiita.com/gx3n-inue/items/4de07dc9e1d90a1cfaa1)

## サンプルの演奏例

### ・「海の見える街」（楽器：シンセ・多声音色１（new age））
```
go run .\Go_PlayBox.go .\PB_sample\PB_uminomierumachi.txt 88
```

### ・「崖の上のポニョ」（楽器：明るい生ピアノ）
```
go run .\Go_PlayBox.go .\PB_sample\PB_ponyo.txt 1
```

第２引数の数値は楽器の指定です。
MIDI音源で定義されている楽器を0～127の値で指定できます。
