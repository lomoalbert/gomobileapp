三，多面手：glVertexAttribPointer 和 glDrawElements
在介绍如何使用 VBO 进行渲染之前，我们先来回顾一下之前使用顶点数组进行渲染用到的函数：

void glVertexAttribPointer (GLuint index, GLint size, GLenum type, GLboolean normalized, GLsizei stride, const GLvoid* ptr);

参数 index ：为顶点数据（如顶点，颜色，法线，纹理或点精灵大小）在着色器程序中的槽位；
参数 size ：指定每一种数据的组成大小，比如顶点由 x, y, z 3个组成部分，纹理由 u, v 2个组成部分；
参数 type ：表示每一个组成部分的数据格式；
参数 normalized ： 表示当数据为法线数据时，是否需要将法线规范化为单位长度，对于其他顶点数据设置为 GL_FALSE 即可。如果法线向量已经为单位长度设置为 GL_FALSE 即可，这样可免去不必要的计算，提升效率；
stride ： 表示上一个数据到下一个数据之间的间隔（同样是以字节为单位），OpenGL ES根据该间隔来从由多个顶点数据混合而成的数据块中跳跃地读取相应的顶点数据；
ptr ：值得注意，这个参数是个多面手。如果没有使用 VBO，它指向 CPU 内存中的顶点数据数组；如果使用 VBO 绑定到 GL_ARRAY_BUFFER，那么它表示该种类型顶点数据在顶点缓存中的起始偏移量。

那 GL_ELEMENT_ARRAY_BUFFER 表示的索引数据呢？那是由以下函数使用的：

void glDrawElements (GLenum mode, GLsizei count, GLenum type, const GLvoid* indices);

参数 mode ：表示描绘的图元类型，如：GL_TRIANGLES，GL_LINES，GL_POINTS；
参数 count ： 表示索引数据的个数；
参数 type ： 表示索引数据的格式，必须是无符号整形值；
indices ：这个参数也是个多面手，如果没有使用 VBO，它指向 CPU 内存中的索引数据数组；如果使用 VBO 绑定到 GL_ELEMENT_ARRAY_BUFFER，那么它表示索引数据在 VBO 中的偏移量。