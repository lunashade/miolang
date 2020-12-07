# miolang

本田未央Advent Calendarの14日目のために作ったみお言語インタープリタです。

## あそびかた

```bash
$ go run . examples/miomio.mio
みおみおみおっ
```

## 中身の話

Whitespaceクローンで、

- SPACE: み
- TAB: お
- LF: っ

に置き換えただけですが、putchar, getcharなどの命令はUTF8のコードポイントを読み書きするようになっています。
