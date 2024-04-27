<template>
    <el-scrollbar height="100%" style="width: 100%;">
        <!-- 标题和搜索框 -->
        <div style="margin-top: 20px; margin-left: 40px; font-size: 2em; font-weight: bold; ">借书证管理
            <el-input v-model="toSearch" :prefix-icon="Search"
                style=" width: 15vw;min-width: 150px; margin-left: 30px; margin-right: 30px; float: right;" clearable />
        </div>

        <!-- 借书证卡片显示区 -->
        <div style="display: flex;flex-wrap: wrap; justify-content: start;">

            <!-- 借书证卡片 -->
            <div class="cardBox" v-for="card in cards" v-show="card.name.includes(toSearch)" :key="card.cardId">
                <div>
                    <!-- 卡片标题 -->
                    <div style="font-size: 25px; font-weight: bold;">No. {{ card.card_id }}</div>

                    <el-divider />

                    <!-- 卡片内容 -->
                    <div style="margin-left: 10px; text-align: start; font-size: 16px;">
                        <p style="padding: 2.5px;"><span style="font-weight: bold;">姓名：</span>{{ card.name }}</p>
                        <p style="padding: 2.5px;overflow: hidden;text-overflow: ellipsis;white-space: nowrap;">
                            <span style="font-weight: bold;">部门：</span>{{ card.department }}</p>
                        <p style="padding: 2.5px;"><span style="font-weight: bold;">类型：</span>{{ card.type }}</p>
                    </div>

                    <el-divider />

                    <!-- 卡片操作 -->
                    <div style="margin-top: 5px;">
                        <el-button type="danger" :icon="Delete" round
                            @click="this.toRemove = card.card_id, this.removeCardVisible = true" >删除</el-button>
                    </div>

                </div>
            </div>

            <!-- 新建借书证卡片 -->
            <el-button class="newCardBox"
                @click="newCardInfo.name = '', newCardInfo.department = '', newCardInfo.type = '学生', newCardVisible = true">
                <el-icon style="height: 50px; width: 50px;">
                    <Plus style="height: 100%; width: 100%;" />
                </el-icon>
            </el-button>

        </div>


        <!-- 新建借书证对话框 -->
        <el-dialog v-model="newCardVisible" title="新建借书证" width="30%" align-center>
            <div style="margin-left: 2vw; font-weight: bold; font-size: 1rem; margin-top: 20px; ">
                姓名：
                <el-input v-model="newCardInfo.name" style="width: 12.5vw;" clearable />
            </div>
            <div style="margin-left: 2vw; font-weight: bold; font-size: 1rem; margin-top: 20px; ">
                部门：
                <el-input v-model="newCardInfo.department" style="width: 12.5vw;" clearable />
            </div>
            <div style="margin-left: 2vw;   font-weight: bold; font-size: 1rem; margin-top: 20px; ">
                类型：
                <el-select v-model="newCardInfo.value" size="middle" style="width: 12.5vw;">
                    <el-option v-for="type in types" :key="type.value" :label="type.label" :value="type.value" />
                </el-select>
            </div>

            <template #footer>
                <span>
                    <el-button @click="newCardVisible = false">取消</el-button>
                    <el-button type="primary" @click="ConfirmNewCard"
                        :disabled="newCardInfo.name.length === 0 || newCardInfo.department.length === 0">确定</el-button>
                </span>
            </template>
        </el-dialog>

        <!-- 删除借书证对话框 -->  
        <el-dialog v-model="removeCardVisible" title="删除借书证" width="30%">
            <span>确定删除<span style="font-weight: bold;">{{ toRemove }}号借书证</span>吗？</span>

            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="removeCardVisible = false">取消</el-button>
                    <el-button type="danger" @click="ConfirmRemoveCard">
                        删除
                    </el-button>
                </span>
            </template>
        </el-dialog>

    </el-scrollbar>
</template>

<script>
import { Delete, Edit, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
    data() {
        return {
            cards: [{ // 借书证列表
                card_id: 1,
                name: '小明',
                department: 'Computer Science',
                type: '学生'
            }, {
                card_id: 2,
                name: '王老师',
                department: '计算机科学与技术学院',
                type: '教师'
            }
            ],
            Delete,
            Search,
            toSearch: '', // 搜索内容
            types: [ // 借书证类型
                {
                    value: 'T',
                    label: '教师',
                },
                {
                    value: 'S',
                    label: '学生',
                }
            ],
            newCardVisible: false, // 新建借书证对话框可见性
            removeCardVisible: false, // 删除借书证对话框可见性
            toRemove: 0, // 待删除借书证号
            newCardInfo: { // 待新建借书证信息
                name: '',
                department: '',
                value: 'S'
            }
        }
    },
    methods: {
        ConfirmNewCard() {
            // 发出POST请求
            axios.post("/card/add",
                { // 请求体
                    name: this.newCardInfo.name,
                    department: this.newCardInfo.department,
                    type: this.newCardInfo.value
                })
                .then(response => {
                    if (response.data.ok) {
                        ElMessage.success("借书证新建成功") // 显示消息提醒
                    } else {
                        ElMessage.error("借书证创建失败: " + response.data.payload.Message) // 显示消息提醒
                    }
                    this.newCardVisible = false // 将对话框设置为不可见
                    this.QueryCards() // 重新查询借书证以刷新页面
                })
        },
        ConfirmRemoveCard() {
            axios.delete("/card/remove",
            { // 请求体
                params: {
                    card_id: this.toRemove
                }
            })
            .then(response => {
                console.log(response.data)
                if (response.data.ok) {
                    ElMessage.success("借书证删除成功") // 显示消息提醒
                } else {
                    ElMessage.error("借书证删除失败: " + response.data.message) // 显示消息提醒
                }
                this.removeCardVisible = false // 将对话框设置为不可见
                this.QueryCards() // 重新查询借书证以刷新页面
            })
        },
        QueryCards() {
            this.cards = [] // 清空列表
            let response = axios.get('/card/query') // 向/card发出GET请求
                .then(response => {
                    if (response.data.ok == false) { // 如果请求失败
                        return ElMessage.error("借书证查询失败: " + response.data.payload.Message) // 显示消息提醒
                    } else {
                        let cards = response.data.payload.cards // 接收响应负载
                        cards.forEach(card => { // 对于每个借书证
                            card.type = (card.type == "T") ? "教师" : "学生" // 将类型转换为中文
                            this.cards.push(card) // 将其加入到列表中
                        })
                    }
                })
        }
    },
    mounted() { // 当页面被渲染时
        //this.QueryCards() // 查询借书证
    }
}

</script>


<style scoped>
.cardBox {
    height: 300px;
    width: 200px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
    text-align: center;
    margin-top: 40px;
    margin-left: 27.5px;
    margin-right: 10px;
    padding: 7.5px;
    padding-right: 10px;
    padding-top: 15px;
}

.newCardBox {
    height: 300px;
    width: 200px;
    margin-top: 40px;
    margin-left: 27.5px;
    margin-right: 10px;
    padding: 7.5px;
    padding-right: 10px;
    padding-top: 15px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
    text-align: center;
}
</style>