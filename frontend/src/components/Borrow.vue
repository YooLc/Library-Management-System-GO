<template>
    <el-scrollbar height="100%" style="width: 100%;">

        <!-- 标题和搜索框 -->
        <div style="margin-top: 20px; margin-left: 40px; font-size: 2em; font-weight: bold;">
            借书记录查询
            <el-input v-model="toSearch" :prefix-icon="Search"
                style=" width: 15vw;min-width: 150px; margin-left: 30px; margin-right: 30px; float: right; ;"
                clearable />
        </div>

        <!-- 查询框 -->
        <div style="width:30%;margin:0 60px; padding-top:5vh;">
            <el-input v-model="this.toQuery" style="display:inline; " placeholder="输入借书证ID"></el-input>
            <el-button style="margin-left: 10px;" type="primary" @click="QueryBorrows">查询</el-button>
        </div>

        <!-- 结果表格 -->
        <el-table v-if="isShow" :data="fitlerTableData" height="600"
            :default-sort="{ prop: 'borrow_time', order: 'ascending' }" :table-layout="'auto'"
            style="width: 100%; margin-left: 50px; margin-top: 30px; margin-right: 50px; max-width: 80vw;">
            <el-table-column prop="card_id" label="借书证 ID" />
            <el-table-column prop="book_id" label="图书 ID" sortable />
            <el-table-column prop="borrow_time" label="借出时间" sortable />
            <el-table-column prop="return_time" label="归还时间" sortable />
            <el-table-column label="操作">
                <template #default="scope">
                    <el-button type="primary" size="small" icon="Plus" :disabled="scope.row.return_time !== '未归还'" @click="
                        returnBookVisible = true; curRow = scope.row;">还书</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!-- 还书确认对话框 -->
        <el-dialog v-model="returnBookVisible" title="归还图书" width="30%">
            <span>确认归还<span style="font-weight: bold;"> {{ curRow.book_id }} 号图书 {{  curRow.title }} </span>吗？</span>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="returnBookVisible = false">取消</el-button>
                    <el-button type="success" @click="ReturnBook()">归还</el-button>
                </span>
            </template>
        </el-dialog>

    </el-scrollbar>
</template>

<script>
import axios from 'axios';
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus';

export default {
    data() {
        return {
            isShow: false, // 结果表格展示状态
            tableData: [{ // 列表项
                card_id: 1,
                book_id: 1,
                borrow_time: 0,
                return_time: 0
            }],
            toQuery: 1, // 待查询内容(对某一借书证号进行查询)
            toSearch: '', // 待搜索内容(对查询到的结果进行搜索)
            returnBookVisible: false, // 还书对话框显示状态
            curRow: null, // 当前行
            Search
        }
    },
    computed: {
        fitlerTableData() { // 搜索规则
            console.log(this.tableData.filter(
                (tuple) =>
                    (this.toSearch == '') || // 搜索框为空，即不搜索
                    tuple.book_id == this.toSearch || // 图书号与搜索要求一致
                    tuple.borrow_time.toString().includes(this.toSearch) || // 借出时间包含搜索要求
                    tuple.return_time.toString().includes(this.toSearch) // 归还时间包含搜索要求
            ))
            return this.tableData.filter(
                (tuple) =>
                    (this.toSearch == '') || // 搜索框为空，即不搜索
                    tuple.book_id == this.toSearch || // 图书号与搜索要求一致
                    tuple.borrow_time.toString().includes(this.toSearch) || // 借出时间包含搜索要求
                    tuple.return_time.toString().includes(this.toSearch) // 归还时间包含搜索要求
            )
        }
    },
    methods: {
        async QueryBorrows() {
            this.tableData = [] // 清空列表
            let response = await axios.get('/borrow/query', { params: { card_id: this.toQuery } }) // 向/borrow发出GET请求，参数为cardID=this.toQuery
            if (response.data.ok == false) {
                ElMessage.error(response.data.message) // 若请求失败，弹出错误信息
                this.isShow = false // 不显示结果列表
                return
            }
            let borrows = response.data.payload.items // 获取响应负载
            if (response.data.payload.count > 0) {
                borrows.forEach(borrow => { // 对于每一个借书记录
                    borrow.borrow_time = new Date(borrow.borrow_time).toLocaleString()
                    if (borrow.return_time != 0) // 若归还时间不为空
                        borrow.return_time = new Date(borrow.return_time).toLocaleString()
                    else
                        borrow.return_time = "未归还"
                    this.tableData.push(borrow) // 将它加入到列表项中
                });
            }
            this.isShow = true // 显示结果列表
        },
        async ReturnBook() {
            let response = await axios.post('/borrow/return', { 
                card_id: this.curRow.card_id,
                book_id: this.curRow.book_id,
                borrow_time: 0,
                return_time: new Date().getTime()
            })
            if (response.data.ok == false) {
                ElMessage.error(response.data.message)
                this.returnBookVisible = false
                return
            }
            ElMessage.success("还书成功")
            this.returnBookVisible = false
            this.QueryBorrows()
        }
    }
}
</script>