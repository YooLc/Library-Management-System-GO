<!-- TODO: YOUR CODE HERE -->
<template>
    <el-scrollbar height="100%" style="width: 100%; height: 100%; ">
        <div style="margin-top: 20px; margin-left: 40px; font-size: 2em; font-weight: bold; ">图书管理</div>
    
        <el-table
        :data="tableData"
        :default-sort="{ prop: 'book_id', order: 'ascending' }"
        style="width: auto; margin-right: 60px; margin-left: 60px;"
        >
            <el-table-column prop="book_id" label="图书编号" sortable/>
            <el-table-column prop="category" label="分类" sortable/>
            <el-table-column prop="title" label="标题" sortable />
            <el-table-column prop="press" label="出版社" sortable/>
            <el-table-column prop="publish_year" label="出版年份" sortable/>
            <el-table-column prop="author" label="作者" sortable/>
            <el-table-column prop="price" label="价格" sortable/>
            <el-table-column prop="stock" label="库存" sortable/>
            <el-table-column label="操作" width="220">
                <template #header>
                    <el-input v-model="search" size="small" placeholder="Type to search" />
                </template>
                <template #default="scope">
                    <el-button type="primary" size="small" icon="Plus" @click="
                        incStockVisible = true; curRow = scope.row;"></el-button>
                    <el-button type="danger" size="small" icon="Delete" @click="
                        removeBookVisible = true; curRow = scope.row;"></el-button>
                    <el-button type="success" size="small" icon="Edit" @click="
                        "></el-button>
                </template>
            </el-table-column>
        </el-table>
    
        <!-- 删除确认对话框 -->
        <el-dialog v-model="removeBookVisible" title="删除图书" width="30%">
            <span>确定删除<span style="font-weight: bold;"> {{ curRow.book_id }} 号图书 {{  curRow.title }} </span>吗？</span>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="removeBookVisible = false">取消</el-button>
                    <el-button type="danger" @click="RemoveBook(curRow.book_id)">删除</el-button>
                </span>
            </template>
        </el-dialog>

        <!-- 修改库存对话框 -->
        <el-dialog v-model="incStockVisible" title="修改库存" width="30%">
            <span>修改<span style="font-weight: bold;"> {{ curRow.book_id }} 号图书 {{  curRow.title }} </span>库存</span>
            <el-input v-model="incStock" style="width: 80%; margin-top: 20px;" placeholder="请输入库存改变量"></el-input>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="removeBookVisible = false">取消</el-button>
                    <el-button type="success" @click="IncBookStock(curRow, incStock)">增加库存</el-button>
                </span>
            </template>
        </el-dialog>
    </el-scrollbar>
</template>

<script lang="ts" setup>
import axios from 'axios'
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'

interface Book {
    book_id: string
    category: string
    title: string
    press: string
    publish_year: number
    author: string
    price: number
    stock: number
}

interface BookQueryCondition {
    category: string | null
    title: string | null
    press: string | null
    minPublishYear: number | null
    maxPublishYear: number | null
    author : string | null
    minPrice: number | null
    maxPrice: number | null
    sortBy: string | null
    sortOrder: string | null
}

const tableData = ref([
    {book_id: '1', category: '计算机', title: '计算机网络', press: '清华大学出版社', author: '谢希仁', price: '5.0', stock: '100'},
    {book_id: '2', category: '计算机', title: '计算机组成原理', press: '清华大学出版社', author: '唐朔飞', price: '6.0', stock: '100'},
    {book_id: '3', category: '计算机', title: '计算机操作系统', press: '清华大学出版社', author: '汤小丹', price: '7.0', stock: '100'},
    {book_id: '4', category: '计算机', title: '计算机图形学', press: '清华大学出版社', author: '王立新', price: '8.0', stock: '100'},
    {book_id: '5', category: '计算机', title: '计算机体系结构', press: '浙江大学出版社', author: '王爱民', price: '9.0', stock: '100'},
    {book_id: '6', category: '计算机', title: '计算机网络', press: '清华大学出版社', author: '谢希仁', price: '5.0', stock: '100'},
])
const nullCondition: BookQueryCondition = {
    category: null,
    title: null,
    press: null,
    minPublishYear: null,
    maxPublishYear: null,
    author: null,
    minPrice: null,
    maxPrice: null,
    sortBy: null,
    sortOrder: null
}

const QueryBooks = (condition: BookQueryCondition) => {
    axios.get('/book/query', { params: condition })
    .then((response) => {
        if (response.data.ok) {
            tableData.value = response.data.payload.results
        } else {
            ElMessage.error('书籍查询失败: ' + response.data.message)
        }
    })
    .catch((error) => {
        ElMessage.error('书籍查询失败: ' + error)
    })
}
onMounted(() => {
    QueryBooks(nullCondition)
})

const IncBookStock = (book, delta_stock) => {
    delta_stock = parseInt(delta_stock)
    if (book.stock + delta_stock < 0) {
        ElMessage.error('修改后库存小于 0, 请确认修改量是否正确')
        return
    }
    axios.put('/book/stock', {book_id: book.book_id, delta_stock: delta_stock})
    .then((res) => {
        if (!res.data.ok) {
            ElMessage.error('库存修改失败: ' + res.data.message)
            return
        }
        ElMessage.success('库存修改成功! 书籍: ' + book.title + ', 当前库存: ' + (book.stock + delta_stock))
        tableData.value = tableData.value.map((b) => {
            if (b.book_id == book.book_id) {
                b.stock += delta_stock
            }
            return b
        })
        incStockVisible.value = false
    })
}

const RemoveBook = (id) => {
    axios.delete('/book/remove', {params: {book_id: id}})
    .then((response) => {
        if (!response.data.ok) {
            ElMessage.error('删除失败: ' + response.data.message)
            return
        }
        ElMessage.success('删除成功')
        tableData.value = tableData.value.filter((book) => book.book_id != id)
        removeBookVisible.value = false
    })
}

const EditBook = () => {
    axios.put('/book/modify', {book_id: '1', category: '计算机', title: '计算机网络', press: '清华大学出版社', author: '谢希仁', price: '5.0', stock: '100'})
    .then((res) => {
        console.log(res.data)
    })
}

const search = ref('')
const removeBookVisible = ref(false)
const incStockVisible = ref(false)
const curRow = ref(null)
const incStock = ref(0)
</script>