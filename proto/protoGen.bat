@ echo off
cd ./protoc/bin
set protoPath="../../"
set clientPath="D:\white_snake\douyin_prj\Assets\Script\GameLogic\Proto"
set serverPath="../../../message"
protoc.exe -I=%protoPath% --csharp_out=%clientPath% --go_out=%serverPath% %protoPath%/*.proto
clang-format.exe -i -style="{AlignConsecutiveAssignments: true,AlignConsecutiveDeclarations: true,AllowShortFunctionsOnASingleLine: None,BreakBeforeBraces: GNU,ColumnLimit: 0,IndentWidth: 4,Language: Proto}" %protoPath%/*.proto
pause