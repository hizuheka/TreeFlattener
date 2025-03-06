package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type fileinfo struct {
	path, date string
}

func (fi *fileinfo) String(outputtype string) string {
	switch outputtype {
	case "1":
		return fmt.Sprintf("%s,%s", fi.path, fi.date)
	case "2":
		return fmt.Sprintf("%s\t%s", fi.path, fi.date)
	case "3":
		return fmt.Sprintf("%s", fi.path)
	default:
		return fmt.Sprintf("%s,%s", fi.path, fi.date)
	}
}

func main() {
	// コマンドライン引数のチェック
	if len(os.Args) < 3 {
		fmt.Println("Usage: program <input_file> <output_file> [<output_type>]")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	outputType := ""
	if len(os.Args) == 4 {
		outputType = os.Args[3]
	}

	// 入力ファイルを開く
	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer input.Close()

	// 出力ファイルを作成
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close()

	// ファイルを1行ずつ読み込むためのスキャナーを作成
	scanner := bufio.NewScanner(input)
	var pathStack []string // 現在のディレクトリパスを保持するスタック

	// ファイル情報を抽出するための正規表現パターン
	filePattern1 := regexp.MustCompile(`^(.+?)\s+(\d{1,3}(?:,\d{3})*)\s+(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})$`) // ファイルサイズあり
	// filePattern2 := regexp.MustCompile(`^(.+?)\s+(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})$`)                        // ファイルサイズなし
	filePattern2 := regexp.MustCompile(`^(.+?)\s+(\d{4}/\d{2}/\d{2})[^\d]+(\d{2}:\d{2}:\d{2})$`) // ファイルサイズなし

	// 出力用のバッファライターを作成
	o := bufio.NewWriter(output)
	defer o.Flush()

	fileinfos := []fileinfo{}

	// 入力ファイルを1行ずつ解析
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue // 空行はスキップ
		}

		// インデントレベルの計算
		indentLevel := calcIndentLevel(line)

		// ルートディレクトリ（最上位のフォルダ）を設定
		if indentLevel == 0 {
			pathStack = []string{getFolderName(line)}
			continue
		}

		// フォルダ、ファイルによって処理を分ける
		if match1 := filePattern1.FindStringSubmatch(line); match1 != nil {
			// ファイルで、サイズありの場合
			filename := getFileName(match1[1]) // ファイル名
			// size := strings.ReplaceAll(match[2], ",", "") // ファイルサイズ（カンマを削除）
			date := match1[3] // 更新日時
			path := strings.Join(pathStack[:indentLevel-1], "\\") + "\\" + filename
			fileinfos = append(fileinfos, fileinfo{path, date})
		} else if match2 := filePattern2.FindStringSubmatch(line); match2 != nil {
			// ファイルで、サイズなしの場合
			filename := getFileName(match2[1])  // ファイル名
			date := match2[2] + " " + match2[3] // 更新日時
			path := strings.Join(pathStack[:indentLevel-1], "\\") + "\\" + filename
			fileinfos = append(fileinfos, fileinfo{path, date})
		} else {
			// フォルダの場合
			name := getFolderName(line)
			// インデントレベルに応じてpathStackを更新
			for len(pathStack) >= indentLevel {
				pathStack = pathStack[:len(pathStack)-1]
			}
			pathStack = append(pathStack, name)
		}
	}

	for _, v := range fileinfos {
		o.WriteString(fmt.Sprintf("%s\r\n", v.String(outputType)))
	}
}

func calcIndentLevel(line string) int {
	replacer := strings.NewReplacer("└", "  ", "│", "  ", "├", "  ")

	// 記号（罫線）を、半角スペース2個に置換
	replaced := replacer.Replace((line))
	// 記号置換後の文字列から、先頭のスペースを削除
	trimmed := strings.TrimLeft(replaced, " ")
	// 1階層でインデント3つなので、3で割る
	return (len(replaced) - len(trimmed)) / 3
}

func getFolderName(line string) string {
	replacer := strings.NewReplacer("└", "", "│", "", "├", "", "< Folder >", "")

	// フォルダ名以外をを削除
	replaced := replacer.Replace((line))
	// スペースをトリム
	return strings.TrimSpace(replaced)
}

func getFileName(s string) string {
	replacer := strings.NewReplacer("└", "", "│", "", "├", "")

	// ファイル名以外をを削除
	replaced := replacer.Replace((s))
	// スペースをトリム
	return strings.TrimSpace(replaced)
}
