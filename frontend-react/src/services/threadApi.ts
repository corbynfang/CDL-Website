import type { AxiosResponse } from 'axios';
import api from './api';
import type { ThreadResponse, ThreadPost } from '../types';

export const threadApi = {
  getThread: async (matchId: number, page = 1, limit = 25): Promise<ThreadResponse> => {
    const response: AxiosResponse<ThreadResponse> = await api.get(
      `/matches/${matchId}/thread`,
      { params: { page, limit } }
    );
    return response.data;
  },

  createPost: async (matchId: number, body: string): Promise<ThreadPost> => {
    const response: AxiosResponse<ThreadPost> = await api.post(
      `/matches/${matchId}/thread/posts`,
      { body }
    );
    return response.data;
  },

  editPost: async (postId: number, body: string): Promise<void> => {
    await api.put(`/thread/posts/${postId}`, { body });
  },

  deletePost: async (postId: number): Promise<void> => {
    await api.delete(`/thread/posts/${postId}`);
  },
};
