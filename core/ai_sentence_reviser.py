#!/usr/bin/env python3
"""
AI味句子识别和改写器
识别并优化写作类action产出中的AI味句子
"""

import re
from typing import List, Tuple, Optional
import logging
from prompts.action_prompts import revision_prompt
from model_base import call_llm  # 使用model_base的call_llm来显示面板

logger = logging.getLogger(__name__)


class AISentenceReviser:
    """AI味句子识别和改写器"""
    
    def __init__(self):
        # AI味句子的正则表达式模式
        self.ai_patterns = [
            # "不是...而是"模式 - 一句话中同时包含"不是"和"而是"
            r'(?=.*不是)(?=.*而是)[^。！？]*[。！？]',
            # "不...而..."模式 - 一句话中同时包含"不"和"而"
            r'(?=.*不)(?=.*而)[^。！？]*[。！？]',
            # "本身就"模式
            r'[^。！？]*?本身就[^。！？]*?[。！？]',
            # "本质"模式
            r'[^。！？]*?本质[^。！？]*?[。！？]',
            # "恰恰"模式
            r'[^。！？]*?恰恰[^。！？]*?[。！？]',
            # "并非...而是"模式 - 一句话中同时包含"并非"和"而是"
            r'(?=.*并非)(?=.*而是)[^。！？]*[。！？]',
            # "不再...而是"模式 - 一句话中同时包含"不再"和"而是"
            r'(?=.*不再)(?=.*而是)[^。！？]*[。！？]',
            # "与其...不如"模式 - 一句话中同时包含"与其"和"不如"
            r'(?=.*与其)(?=.*不如)[^。！？]*[。！？]',
            # "那...好像"模式 - 一句话中同时包含"那"和"好像"
            r'(?=.*那)(?=.*好像)[^。！？]*[。！？]',
            # "后来我发现"模式
            r'[^。！？]*?后来我发现[^。！？]*?[。！？]',
        ]
    
    def identify_ai_sentences(self, text: str) -> List[str]:
        """
        识别文本中的AI味句子
        
        Args:
            text: 要分析的文本
            
        Returns:
            AI味句子列表
        """
        ai_sentences = []
        
        for pattern in self.ai_patterns:
            matches = re.findall(pattern, text)
            for match in matches:
                # 清理句子，去掉多余的空白
                sentence = match.strip()
                if sentence and sentence not in ai_sentences:
                    ai_sentences.append(sentence)
                    logger.debug(f"识别到AI味句子: {sentence[:30]}...")
        
        return ai_sentences
    
    def extract_context_for_sentence(self, text: str, target_sentence: str) -> str:
        """
        提取句子的上下文（前后2句）
        
        Args:
            text: 完整文本
            target_sentence: 目标句子
            
        Returns:
            包含上下文的文本片段
        """
        # 按句子分割文本
        sentences = re.split(r'[。！？]', text)
        sentences = [s.strip() for s in sentences if s.strip()]
        
        # 找到目标句子的位置
        target_index = -1
        for i, sentence in enumerate(sentences):
            if target_sentence.replace('。', '').replace('！', '').replace('？', '').strip() in sentence:
                target_index = i
                break
        
        if target_index == -1:
            logger.warning(f"未找到目标句子: {target_sentence[:30]}...")
            return target_sentence
        
        # 提取前后2句
        start_index = max(0, target_index - 2)
        end_index = min(len(sentences), target_index + 3)
        
        context_sentences = sentences[start_index:end_index]
        
        # 重新组合成文本，并标记目标句子
        context_text = ""
        for i, sentence in enumerate(context_sentences):
            actual_index = start_index + i
            if actual_index == target_index:
                context_text += f"<bad_sentence>{sentence}。</bad_sentence>"
            else:
                context_text += f"{sentence}。"
        
        return context_text
    
    def revise_sentence(self, context_text: str) -> Optional[str]:
        """
        使用LLM改写句子
        
        Args:
            context_text: 包含<bad_sentence>标签的上下文
            
        Returns:
            改写后的句子，如果失败则返回None
        """
        try:
            # 调用LLM进行改写
            response = call_llm(
                system_prompt=revision_prompt,
                user_prompt=context_text,
                temperature=0.7,
                max_output_tokens=2000
            )
            
            # 提取<revised_bad_sentence>中的内容
            match = re.search(r'<revised_bad_sentence>(.*?)</revised_bad_sentence>', response, re.DOTALL)
            if match:
                revised_sentence = match.group(1).strip()
                logger.debug(f"句子改写成功: {revised_sentence[:30]}...")
                return revised_sentence
            else:
                logger.warning("未找到<revised_bad_sentence>标签")
                return None
                
        except Exception as e:
            logger.error(f"句子改写失败: {e}")
            return None
    
    def revise_text(self, text: str) -> str:
        """
        改写文本中的所有AI味句子
        
        Args:
            text: 原始文本
            
        Returns:
            改写后的文本
        """
        if not text or not text.strip():
            return text
        
        # 识别AI味句子
        ai_sentences = self.identify_ai_sentences(text)
        if not ai_sentences:
            logger.debug("未发现AI味句子")
            return text
        
        logger.info(f"发现 {len(ai_sentences)} 个AI味句子，开始改写...")
        
        revised_text = text
        
        # 逐个改写句子
        for ai_sentence in ai_sentences:
            # 提取上下文
            context_text = self.extract_context_for_sentence(text, ai_sentence)
            
            # 改写句子
            revised_sentence = self.revise_sentence(context_text)
            
            if revised_sentence:
                # 替换原句子
                # 移除原句子的标点符号进行匹配
                original_clean = ai_sentence.replace('。', '').replace('！', '').replace('？', '').strip()
                revised_text = revised_text.replace(ai_sentence, revised_sentence)
                logger.info(f"句子改写完成: {original_clean[:20]}... → {revised_sentence[:20]}...")
            else:
                logger.warning(f"句子改写失败，保持原样: {ai_sentence[:30]}...")
        
        return revised_text


# 创建全局实例
_reviser = AISentenceReviser()


def revise_ai_sentences(text: str) -> str:
    """
    改写文本中的AI味句子
    
    Args:
        text: 原始文本
        
    Returns:
        改写后的文本
    """
    return _reviser.revise_text(text)


def should_revise_action_type(action_type: str) -> bool:
    """
    判断是否需要对指定的action类型进行AI味句子改写
    
    Args:
        action_type: action类型
        
    Returns:
        是否需要改写
    """
    writing_actions = ['xhs_post', 'wechat_article', 'tiktok_script']
    return action_type in writing_actions 