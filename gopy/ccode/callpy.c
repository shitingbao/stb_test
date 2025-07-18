// ccode/call_python.c

#include <Python.h>
#include <stdio.h>
#include <string.h>
#include "callpy.h"

void call_python(const char* name, int times) {
    // 初始化Python解释器
    Py_Initialize();

    // 设置Python模块路径
    PyRun_SimpleString("import sys"); 
    PyRun_SimpleString("sys.path.append(\"./pycode\")");

    // 导入模块
    PyObject* pName = PyUnicode_FromString("greettest");
    PyObject* pModule = PyImport_Import(pName);
    Py_DECREF(pName);

    if (!pModule) {
        PyErr_Print();
        fprintf(stderr, "[C] Failed to load \"greettest.py\"\n");
        Py_Finalize();
        return;
    }

    // 获取Greeter类
    PyObject* pClass = PyObject_GetAttrString(pModule, "Greeter");
    if (!pClass || !PyCallable_Check(pClass)) {
        PyErr_Print();
        fprintf(stderr, "[C] Failed to get Greeter class\n");
        Py_DECREF(pModule);
        Py_Finalize();
        return;
    }

    // 创建Greeter对象
    PyObject* pArgs = PyTuple_Pack(1, PyUnicode_FromString(name));
    PyObject* pInstance = PyObject_CallObject(pClass, pArgs);
    Py_DECREF(pArgs);

    if (!pInstance) {
        PyErr_Print();
        fprintf(stderr, "[C] Failed to create Greeter instance\n");
        Py_DECREF(pClass);
        Py_DECREF(pModule);
        Py_Finalize();
        return;
    }

    // 调用greet方法
    PyObject* pValue = PyObject_CallMethod(pInstance, "greet", "i", times);
    if (pValue) {
        const char* result = PyUnicode_AsUTF8(pValue);
        printf("[C] Python returned: %s\n", result);
        Py_DECREF(pValue);
    } else {
        PyErr_Print();
        fprintf(stderr, "[C] Failed to call greet()\n");
    }

    Py_DECREF(pInstance);
    Py_DECREF(pClass);
    Py_DECREF(pModule);
    Py_Finalize();
}
