package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	// コマンドライン引数のチェック
	if len(os.Args) < 3 {
		fmt.Println("Usage: program <input_file> <output_file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

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
	var basePath string    // ルートディレクトリの名前
	var pathStack []string // 現在のディレクトリパスを保持するスタック

	// ファイル情報を抽出するための正規表現パターン
	filePattern := regexp.MustCompile(`^(.+?)\s+(\d{1,3}(?:,\d{3})*)\s+(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})$`)
	// filePattern := regexp.MustCompile(`[└├]\s(.+?)\s+(\d{1,3}(?:,\d{3})*)\s+(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})$`)

	// 出力用のバッファライターを作成
	o := bufio.NewWriter(output)
	defer o.Flush()

	replacer := strings.NewReplacer("└", "  ", "│", "  ", "├", "  ")

	// 入力ファイルを1行ずつ解析
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue // 空行はスキップ
		}

		// 先頭の空白を削除し、インデントレベルを計算
		// trimmed = strings.TrimLeft(line, " ")
		replaced := replacer.Replace((line))
		trimmed := strings.TrimLeft(replaced, " ")
		fmt.Println("trimmed=" + trimmed)
		indentLevel := len(replaced) - len(trimmed)
		fmt.Printf("indentLevel=%d,len(line)=%d,len(trimmed)=%d\n", indentLevel, len(line), len(trimmed))

		// ルートディレクトリ（最上位のフォルダ）を設定
		if indentLevel == 0 {
			trimmed = strings.ReplaceAll(trimmed, "< Folder >", "")
			trimmed = strings.TrimSpace(trimmed)
			basePath = trimmed
			pathStack = []string{basePath}
			continue
		}

		// ファイル情報の行かどうかを判定
		match := filePattern.FindStringSubmatch(trimmed)
		if match != nil {
			filename := match[1]                          // ファイル名
			size := strings.ReplaceAll(match[2], ",", "") // ファイルサイズ（カンマを削除）
			date := match[3]                              // 更新日時
			fullPath := strings.Join(pathStack, "\\") + "\\" + filename
			o.WriteString(fmt.Sprintf("%s,%s,%s\n", fullPath, size, date))
			continue
		}

		// フォルダ名の前にある記号を削除
		trimmed = strings.ReplaceAll(trimmed, "< Folder >", "")
		trimmed = strings.TrimSpace(trimmed)

		// インデントレベルに応じてpathStackを更新
		for len(pathStack) > indentLevel/3+1 {
			pathStack = pathStack[:len(pathStack)-1-1-1]
		}
		pathStack = append(pathStack, trimmed)
		fmt.Println(pathStack)
	}
}
