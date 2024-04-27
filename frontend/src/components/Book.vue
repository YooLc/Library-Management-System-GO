<!-- TODO: YOUR CODE HERE -->
<template>
    <el-scrollbar height="100%" style="width: 100%; height: 100%; ">
        <div style="margin-top: 20px; margin-left: 40px; font-size: 2em; font-weight: bold; ">图书管理</div>

        <!-- 条件搜索区域 -->
        <el-form :inline="true" style="margin-top: 20px; margin-left: 40px;">
            <el-form-item label="分类">
                <el-input v-model="condition.category" placeholder="分类"></el-input>
            </el-form-item>
            <el-form-item label="标题">
                <el-input v-model="condition.title" placeholder="标题"></el-input>
            </el-form-item>
            <el-form-item label="出版社">
                <el-input v-model="condition.press" placeholder="出版社"></el-input>
            </el-form-item>
            <el-form-item label="出版年份">
                <el-input-number size="mini" v-model="condition.min_publish_year" placeholder="最小年份"></el-input-number>
                <p style="color: gray; font-size: 12px; margin: 0 20px;"> ~ </p>
                <el-input-number size="mini" v-model="condition.max_publish_year" placeholder="最大年份"></el-input-number>
            </el-form-item>
            <el-form-item label="作者">
                <el-input v-model="condition.author" placeholder="作者"></el-input>
            </el-form-item>
            <el-form-item label="价格">
                <el-input width="20px" v-model="condition.min_price" placeholder="最小价格" style="width: 100px;"></el-input>
                <p style="color: gray; font-size: 12px; margin: 0 20px;"> ~ </p>
                <el-input size="mini" v-model="condition.max_price" placeholder="最大价格" style="width: 100px;"></el-input>
            </el-form-item>
            <el-form-item label="排序">
                <el-select v-model="condition.sort_by" placeholder="排序" style="width: 100px;">
                    <el-option label="图书编号" value="book_id"></el-option>
                    <el-option label="分类" value="category"></el-option>
                    <el-option label="标题" value="title"></el-option>
                    <el-option label="出版社" value="press"></el-option>
                    <el-option label="出版年份" value="publish_year"></el-option>
                    <el-option label="作者" value="author"></el-option>
                    <el-option label="价格" value="price"></el-option>
                    <el-option label="库存" value="stock"></el-option>
                </el-select>
                <el-select v-model="condition.sort_order" placeholder="排序方式" style="width: 100px; margin-left: 20px;">
                    <el-option label="升序" value="asc"></el-option>
                    <el-option label="降序" value="desc"></el-option>
                </el-select>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="QueryBooks(condition)">搜索</el-button>
            </el-form-item>
        </el-form>
    
        <el-table
        :data="tableData.filter(data => !search || data.title.toLowerCase().includes(search.toLowerCase()))"
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
            <el-table-column label="操作" width="350">
                <template #header>
                    <el-input v-model="search" size="small" placeholder="Type to search" />
                </template>
                <template #default="scope">
                    <el-button type="primary" size="small" icon="Minus" @click="
                        borrowBookVisible = true; curRow = scope.row;">借书</el-button>
                    <el-button type="success" size="small" icon="Plus" @click="
                        incStockVisible = true; curRow = scope.row;">入库</el-button>
                    <el-button type="warning" size="small" icon="Edit" @click="
                        modifyBookVisible = true; curRow = scope.row; ">修改</el-button>
                    <el-button type="danger" size="small" icon="Delete" @click="
                        removeBookVisible = true; curRow = scope.row;">删除</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!-- 借书对话框 -->
        <el-dialog v-model="borrowBookVisible" title="借书" width="30%">
            <span>借书证号: 
                <el-input-number size="small" v-model="borrowCardId" placeholder="借书证号"></el-input-number>
            </span>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="borrowBookVisible = false">取消</el-button>
                    <el-button type="success" @click="BorrowBook(curRow, borrowCardId)">借书</el-button>
                </span>
            </template>
        </el-dialog>
    
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
        <el-dialog v-model="incStockVisible" title="增加库存" width="30%">
            <span>增加<span style="font-weight: bold;"> {{ curRow.book_id }} 号图书 {{  curRow.title }} </span>库存</span>
            <p>请输入库存改变量:</p>
            <el-input-number v-model="incStock" style="width: 30%; margin-top: 20px;" placeholder="请输入库存改变量"></el-input-number>
            
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="incStockVisible = false">取消</el-button>
                    <el-button type="success" @click="IncBookStock(curRow, incStock)">增加库存</el-button>
                </span>
            </template>
        </el-dialog>

        <!-- 修改书籍信息对话框 -->
        <el-dialog v-model="modifyBookVisible" title="编辑书籍信息">
            <el-form :model="curRow" label-width="80px">
                <el-form-item label="分类">
                    <el-input v-model="curRow.category"></el-input>
                </el-form-item>
                <el-form-item label="标题">
                    <el-input v-model="curRow.title"></el-input>
                </el-form-item>
                <el-form-item label="出版社">
                    <el-input v-model="curRow.press"></el-input>
                </el-form-item>
                <el-form-item label="出版年份">
                    <el-input v-model="curRow.publish_year"></el-input>
                </el-form-item>
                <el-form-item label="作者">
                    <el-input v-model="curRow.author"></el-input>
                </el-form-item>
                <el-form-item label="价格">
                    <el-input v-model="curRow.price"></el-input>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="modifyBookVisible = false">取消</el-button>
                    <el-button type="success" @click="Edit_book(curRow)">修改</el-button>
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
    book_id: number
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
    min_publish_year: number | null
    max_publish_year: number | null
    author : string | null
    min_price: number | null
    max_price: number | null
    sort_by: string | null
    sort_order: string | null
}

const tableData = ref([
    {book_id: 1, category: '计算机', title: '计算机网络', press: '清华大学出版社', author: '谢希仁', price: 5.0, stock: 100},
    {book_id: 2, category: '计算机', title: '计算机组成原理', press: '清华大学出版社', author: '唐朔飞', price: 6.0, stock: 100},
    {book_id: 3, category: '计算机', title: '计算机操作系统', press: '清华大学出版社', author: '汤小丹', price: 7.0, stock: 100},
    {book_id: 4, category: '计算机', title: '计算机图形学', press: '清华大学出版社', author: '王立新', price: 8.0, stock: 100},
    {book_id: 5, category: '计算机', title: '计算机体系结构', press: '浙江大学出版社', author: '王爱民', price: 9.0, stock: 100},
    {book_id: 6, category: '计算机', title: '计算机网络', press: '清华大学出版社', author: '谢希仁', price: 5.0, stock: 100},
])
const nullCondition: BookQueryCondition = {
    category: null,
    title: null,
    press: null,
    min_publish_year: null,
    max_publish_year: null,
    author: null,
    min_price: null,
    max_price: null,
    sort_by: null,
    sort_order: null
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
            removeBookVisible.value = false
            return
        }
        ElMessage.success('删除成功')
        tableData.value = tableData.value.filter((book) => book.book_id != id)
        removeBookVisible.value = false
    })
}

const Edit_book = (book) => {
    book.publish_year = parseInt(book.publish_year)
    book.price = parseFloat(book.price)

    axios.put('/book/modify', book)
    .then((res) => {
        if (!res.data.ok) {
            ElMessage.error('修改失败: ' + res.data.message)
            return
        }
        ElMessage.success('修改成功')
        modifyBookVisible.value = false
    })
}

const BorrowBook = (book: Book, card_id) => {
    axios.post('/borrow/add', {
        book_id: book.book_id,
        card_id: card_id,
        borrow_time: new Date().getTime(),
        return_time: 0,
    })
    .then((res) => {
        if (!res.data.ok) {
            ElMessage.error('借书失败: ' + res.data.message)
            borrowBookVisible.value = false
            return
        }
        ElMessage.success('借书成功')
        tableData.value = tableData.value.map((b) => {
            if (b.book_id === book.book_id) {
                b.stock -= 1
            }
            return b
        })
        borrowBookVisible.value = false
    })
}

const search = ref('')
const removeBookVisible = ref(false)
const incStockVisible = ref(false)
const curRow = ref(null)
const incStock = ref(0)
const modifyBookVisible = ref(false)
const condition = ref(nullCondition)
const borrowCardId = ref(1)
const borrowBookVisible = ref(false)
</script>