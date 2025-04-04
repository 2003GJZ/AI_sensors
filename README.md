# 老旧设备AI监控识别系统

## 项目背景  
随着工业设备的老化，许多老旧仪器（如压力表、水表、温度表等）由于缺乏电气化功能，依赖人力维护，成本较高且效率低下。为了解决这一问题，本项目通过 **IoT 设备图像采集** 和 **AI 识别技术**，实现对老旧设备的低成本升级与智能化管理。

---

## 系统架构  

本项目分为两种主要数据采集和处理方式：  

1. **图像采集与AI识别流程**  
   - IoT 设备采集仪器表盘图像，通过 **FTP** 协议上传图片至后端服务。  
   - 上传完成后，通过 **HTTP 通知** 后端服务。  
   - 后端服务通过 **NFS 文件共享** 与 AI 服务器进行联动，将图片路径发送至 AI 服务器。  
   - AI 服务器识别后，将结果返回给后端服务。

2. **设备参数报文上传流程**  
   - 某些IoT设备无需上传图片，而是通过 **MQTT 协议** 上报设备参数（例如电表数据）。  
   - 数据以 **DLT645-2007 报文** 格式传输，后端服务解析报文并提取数据后完成上报处理。

---

## 功能模块  

### 图像处理模块  
- **FTP 图片接收**：接收来自 IoT 设备的图片。  
- **文件共享管理**：通过 NFS 实现后端服务与 AI 服务器共享图片文件。  
- **AI 识别结果处理**：接收 AI 服务器返回的识别结果并存储。

### 参数解析模块  
- **MQTT 消息接收**：接收 IoT 设备上传的设备参数。  
- **DLT645-2007 报文解析**：解析电表等设备上传的报文，提取关键参数。  
- **数据上报**：将解析后的数据进行后续存储或通知。  

---

## 技术栈  
- **后端开发语言**：Go  
- **通信协议**：  
  - FTP (文件传输)  
  - HTTP (通知)  
  - MQTT (参数上报)  
- **文件共享**：NFS (Network File System)  
- **报文解析**：DLT645-2007 协议解析  
- **AI 模型集成**：通过 AI 服务器完成数据识别与返回  

---

## 部署说明  

### 环境要求  
- **操作系统**：Linux (推荐 Ubuntu 24)  
- **必要软件**：  
  - FTP 服务  
  - NFS 服务  
  - MQTT Broker (如 Mosquitto)  
- **AI 服务器**：支持自定义识别模型并提供 API 接口  

### 配置步骤  
1. **FTP 服务**：配置用于接收 IoT 设备上传的图片。  
2. **NFS 文件共享**：确保后端服务与 AI 服务器共享同一目录。  
3. **后端服务启动**：运行后端服务以处理图片与报文数据。  
4. **AI 服务对接**：确保 AI 服务器 API 可正常识别图片并返回结果。  

---

## 数据流示意图  

```mermaid
graph TD
    IoT设备1 -->|FTP 上传图片| 后端服务
    IoT设备1 -->|HTTP 通知| 后端服务
    后端服务 -->|NFS 文件路径| AI服务器
    AI服务器 -->|识别结果返回| 后端服务

    IoT设备2 -->|MQTT 报文| 后端服务
    后端服务 -->|解析结果| 数据存储
