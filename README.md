# 标记文档
通过标记方式,收集数据,再通过模版选择生成新的文档,相比于传统的应用方案
优势:
1. 实现单一数据源原则,能确保数据一致性、和准确性
2. 方便修改,在数据源修改后,所有引用的地方都能自动修改
3. 适应性强,对提取数据的原始文档没有特别要求,能很好的利用markdown、text等文档,方便和产品文档、技术文档结合
4. 文档格式灵活,能适应各种格式文档
劣势:
1. 需要记住一些语法
2. 数据不能及时展示,需要解析后才能形成最终效果
3. 数据嵌套问题,部分数据即来自上次的模板渲染后,又成为下一次数据提取的来源,并且这种嵌套混合在一个文档中,对设计、实现是一个巨大挑战