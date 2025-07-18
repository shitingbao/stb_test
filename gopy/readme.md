golang 调用 python类方法

golang 调用c 方法 ，然后 c 的库中调用 python 方法

1.golang调用c方法
```
  add.c + add.h  →(编译)→  add.o  →(打包)→  libadd.a
     ↑             ↑         ↓         ↓
  C源代码     编译为目标文件     静态链接库
     ↓             ↓         ↑         ↑
main.go     通过cgo调用它们 ←←←←←←←←←←

```

2.详细解释一下 .c、.h、.o、.a 这四种文件的含义和作用

.c —— C 源代码文件
```
就是实现函数逻辑的地方。

比如你写的 int add(int a, int b) { return a + b; } 就是在 .c 文件中。

相当于 Python 的 .py、Go 的 .go。

文件例子：add.c
内容：真正的函数实现。
```
📄 .h —— C 头文件（Header）
```
就是声明接口，告诉别人你有哪些函数。
```
.h 是 .c 的“对外说明书”，函数原型放在里面。
```
你在 Go 中的 #include "add.h" 就是用这个文件来告诉 Go 侧：“有一个 C 函数叫 add(int, int)”。

文件例子：add.h
内容：声明 int add(int, int);。
```
🧱 .o —— 目标文件（Object File）
```
.c 文件被编译器编译后的中间产物。

包含机器码，但还不能独立运行。

是下一步打包成 .a 的基础。

文件例子：add.o
命令：gcc -c add.c -o add.o
含义：将 add.c 编译成 add.o，不进行链接。
```

📦 .a —— 静态库（Archive）
```
把多个 .o 文件打包成一个可复用的库文件。

.a 可以在不同的程序中复用，就像 Python 的 .pyc + .zip 打包。

是给 Go 链接器用的，Go 运行时通过 libadd.a 找到并调用 add.o 中的 add()。

文件例子：libadd.a
命令：ar rcs libadd.a add.o
含义：把 add.o 打包成 libadd.a，供链接使用。

```

🚀 总结一张表
```
后缀	文件类型	内容/用途	相当于
.c	C 源文件	写函数逻辑（实现）	Python 的 .py
.h	头文件	写函数声明（接口）	Python 的接口说明或 .pyi
.o	目标文件	编译后的机器码文件，不能单独执行	.pyc（但未封装）
.a	静态链接库	把 .o 打包为库，供别人调用	.zip 或打包好的模块
```



特殊注释里面含义解释

```
 #cgo CFLAGS: -I.
这表示给 C 编译器传入编译参数：

-I.：表示 包含当前目录（.）作为头文件搜索路径。

等价于 C 里写的：gcc -I.

也就是说，当 #include "add.h" 时，它会去当前目录找 add.h。

🔹 #cgo LDFLAGS: -L. -ladd
这行给 链接器传参数：

-L.：表示 链接器从当前目录查找库文件（.a 或 .so）

-ladd：表示要链接一个库叫 libadd.a 或 libadd.so

📌 注意：-ladd 实际链接的是：

libadd.a（静态库），或

libadd.so（动态库）

-l<name> 查找的是 lib<name>.{a,so}，这就是 Linux/C 的命名规范。
```

# 生成代码过程
在你的c代码目录执行，这里是在 ccode 目录下执行（golang 代码中特殊注释只是引用，不需要改）：
```
#静态链接
gcc -c add.c -o add.o
ar rcs libadd.a add.o

# 有py的代码第一句生成o文件（中间文件）
gcc -c -I/Library/Frameworks/Python.framework/Versions/3.11/include/python3.11 -o ccode/callpy.o  ccode/callpy.c
# 第二句生成链接文件（a文件）
ar rcs ccode/libcallpy.a ccode/callpy.o

#动态链接
gcc -fPIC -c add.c -o add.o
gcc -shared -o libadd.so add.o
```

会生成对应后缀为 “o” 的文件以及 “a” 的文件

注意：
1.如果修改了 C 代码，需要将go编译出来的文件删除，重新生成，因为go代码不变，不会去重新生成可执行文件，同时，需要重新生成 .o 和 .a 文件
2.即使编译好 go文件，也需要py的代码和基本结构，因为go代码里面的引用，而C相关的都不需要
最终需要的文件
（go生成的可执行文件+完整的结构的python文件即可）
在这个目录下就是（gopy执行文件+pycode/greettest.py文件结构和文件）即可