package common

//官方说明
// 为什么Go没有variant（变种）类型？
// variant类型，也称作代数类型，提供一种方式来指定一个值可能是一些类型集中的一种，但也只能是这些类型。系统编程中一个常见例子是指定错误，例如，网络错误、安全错误或应用错误，并允许调用者通过检查错误类型来区分问题的来源。另一个例子是语义树中的每一个节点都可以是不同的类型：声明、陈述、赋值等等。

// 我们考虑过为Go添加variant类型，但经过一番讨论之后决定把它们移除掉，因为它们与接口的重叠会令人困惑。想一想如果一个variant类型的元素是它们自己的接口会发生什么？

// 此外，该语言的已经涵盖了variant类型的一些目的。上面错误类型的例子很容易表达使用接口值来报错错误并用类型来区分场景。语义树的例子也可行，尽管没有这么优雅。

//Variant 类借鉴于Qt的QVariant类，类似于Boost的any类。它把常用类型使用一个类包装起来，这样使用QVector等容器时，其内部就可以存储不同的数据
type Variant struct{}

//StaticVariant ..
type StaticVariant struct{}

//其实就是一些message的集合，如果做成接口类？
// C++用法
// using net_message = static_variant<handshake_message,
// chain_size_message,
// go_away_message,
// time_message,
// notice_message,
// request_message,
// sync_request_message,
// signed_block,         // which = 7
// packed_transaction>;  // which = 8
