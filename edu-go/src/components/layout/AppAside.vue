<script setup lang="ts">
import { User, Setting } from '@element-plus/icons-vue';

</script>


<template>
    <el-aside>
        <el-menu router>
            <!-- 设置子菜单 -->
            <el-sub-menu>
                <!-- 由于本身这个菜单点开只显示子菜单 所以这里需要特殊处理 -->
                <template #title>
                    <el-icon>
                        <Setting />
                    </el-icon> <span>权限管理</span>
                    <!-- 权限管理的子菜单项依旧使用el-menu-item -->
                    <!-- 使用router index来使得点击之后跳转 -->
                    <el-menu-item index="/menus">
                        <el-icon>
                            <Setting />
                        </el-icon> <span>菜单列表</span>
                    </el-menu-item>
                    <el-menu-item>
                        <el-icon>
                            <Setting />
                        </el-icon> <span>菜单列表</span>
                    </el-menu-item>
                    <el-menu-item index="resources">
                        <el-icon>
                            <Setting />
                        </el-icon> <span>资源列表</span>
                    </el-menu-item>
                    <el-menu-item index="roles">
                        <el-icon>
                            <Setting />
                        </el-icon> <span>角色列表</span>
                    </el-menu-item>
                </template>
            </el-sub-menu>
            <!-- <el-menu-item>
                <el-icon>
                    <Setting />
                </el-icon> <span>权限管理</span>
            </el-menu-item> -->

            <el-menu-item>
                <el-icon>
                    <User />
                </el-icon> <span>用户管理</span>
            </el-menu-item>

            <el-menu-item>
                <el-icon>

                </el-icon> <span>课程管理</span>
            </el-menu-item>

            <el-sub-menu>
                <template #title>
                    <el-icon>

                    </el-icon><span>广告管理</span>
                </template>
                <el-menu-item>
                    <el-icon>
                        <Setting />
                    </el-icon> <span>广告管理</span>
                </el-menu-item>
                <el-menu-item>
                    <el-icon>
                        <Setting />
                    </el-icon> <span>广告列表</span>
                </el-menu-item>
            </el-sub-menu>

        </el-menu>
    </el-aside>
</template>

<style lang="scss" scoped>
.el-aside {
    background-color: #eee;
    height: 100vh;
    width: auto;
}

.el-menu {
    background-color: #eee;
    width: 200px;
}
</style>