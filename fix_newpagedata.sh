#!/bin/bash

# NewPageData 替换为 NewPageData2 主/子菜单索引组合
declare -A replacements=(
  ["0"]="0"
  ["1"]="0"
  ["2"]="2"
  ["3"]="2"
  ["4"]="2"
  ["5"]="5"
  ["6"]="5"
  ["7"]="7"
  ["8"]="7"
  ["9"]="7"
  ["10"]="7"
  ["11"]="7"
  ["12"]="7"
  ["13"]="7"
  ["14"]="7"
  ["15"]="7"
  ["16"]="7"
  ["17"]="7"
  ["18"]="7"
  ["19"]="7"
  ["20"]="7"
  ["21"]="7"
  ["22"]="22"
  ["23"]="22"
  ["24"]="22"
  ["25"]="22"
)

# 遍历所有 .go 文件并替换匹配调用
echo "🚀 正在替换 NewPageData(...) 为 NewPageData2(...)"

for gofile in $(find vul -type f -name "*.go"); do
  for sub in "${!replacements[@]}"; do
    main="${replacements[$sub]}"
    sed -i -E "s/NewPageData\\($sub,/NewPageData2($main, $sub,/g" "$gofile"
  done
done

echo "✅ 替换完成！请运行 go run . 验证页面高亮和菜单展开效果"
