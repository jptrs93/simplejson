# sjson

`sjson` is a simple light weight package for parsing and manipulating json data of unknown structure. Once parsed, the structure of the json is captured in a general purpose struct allowing it to be repeatedly queried and manipulated efficiently. The parsing is performed in a single pass of the bytes and may be done re-using a byte slice without copying. 