import { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { useAuth } from "../../context/AuthContext";
import { threadApi } from "../../services/threadApi";
import type { ThreadPost } from "../../types";

interface Props {
	matchId: number;
}

export default function MatchThread({ matchId }: Props) {
	const { user, session, needsProfileSetup, openAuthModal } = useAuth();
	const [posts, setPosts] = useState<ThreadPost[]>([]);
	const [total, setTotal] = useState(0);
	const [page, setPage] = useState(1);
	const [body, setBody] = useState("");
	const [editingId, setEditingId] = useState<number | null>(null);
	const [editBody, setEditBody] = useState("");
	const [error, setError] = useState<string | null>(null);
	const [loading, setLoading] = useState(false);
	const limit = 25;

	const fetchPosts = useCallback(async (targetPage = page) => {
		try {
			const data = await threadApi.getThread(matchId, targetPage, limit);
			setPosts(data.data);
			setTotal(data.pagination.total);
		} catch {
			// non-blocking
		}
	}, [matchId, page]);

	useEffect(() => {
		fetchPosts();
	}, [fetchPosts]);

	async function handleSubmit(e: React.FormEvent) {
		e.preventDefault();
		if (!body.trim()) return;
		setError(null);
		setLoading(true);
		try {
			await threadApi.createPost(matchId, body.trim());
			setBody("");
			setPage(1);
			await fetchPosts(1);
		} catch (err: unknown) {
			const msg = axios.isAxiosError<{ error: string }>(err)
				? (err.response?.data?.error ?? "Failed to post")
				: "Failed to post";
			setError(msg);
		}
		setLoading(false);
	}

	async function handleEdit(postId: number) {
		if (!editBody.trim()) return;
		try {
			await threadApi.editPost(postId, editBody.trim());
			setEditingId(null);
			await fetchPosts();
		} catch (err: unknown) {
			const msg = axios.isAxiosError<{ error: string }>(err)
				? (err.response?.data?.error ?? "Failed to edit")
				: "Failed to edit";
			setError(msg);
		}
	}

	async function handleDelete(postId: number) {
		if (!confirm("Delete this post?")) return;
		try {
			await threadApi.deletePost(postId);
			await fetchPosts();
		} catch (err: unknown) {
			const msg = axios.isAxiosError<{ error: string }>(err)
				? (err.response?.data?.error ?? "Failed to delete")
				: "Failed to delete";
			setError(msg);
		}
	}

	const totalPages = Math.ceil(total / limit);
	const isOwner = (post: ThreadPost) => user?.id === post.user.supabase_uid;

	return (
		<div className="mt-8">
			{/* Section header */}
			<p className="text-xs uppercase tracking-widest text-[#737373] mb-4">
				Discussion · {total} comment{total !== 1 ? "s" : ""}
			</p>

			<div className="border border-[#1a1a1a] bg-[#111111]">
				{/* Compose / sign-in prompt */}
				<div className="border-b border-[#1a1a1a] p-4">
					{session && !needsProfileSetup ? (
						<form onSubmit={handleSubmit}>
							<textarea
								placeholder="Write a comment…"
								value={body}
								onChange={(e) => setBody(e.target.value)}
								maxLength={2000}
								required
								rows={3}
								className="w-full bg-[#0a0a0a] border border-[#1a1a1a] text-white text-sm placeholder-[#404040] px-3 py-2 resize-none focus:outline-none focus:border-[#404040] transition-colors"
							/>
							<div className="flex items-center justify-between mt-2">
								<span className="text-[#404040] text-xs font-mono">
									{body.length}/2000
								</span>
								<button
									type="submit"
									disabled={loading || !body.trim()}
									className="px-4 py-1.5 text-xs font-grotesk font-semibold uppercase tracking-widest bg-white text-black disabled:opacity-30 disabled:cursor-not-allowed hover:bg-[#e0e0e0] transition-colors"
								>
									{loading ? "Posting…" : "Post"}
								</button>
							</div>
							{error && <p className="text-red-400 text-xs mt-2">{error}</p>}
						</form>
					) : session && needsProfileSetup ? (
						<button
							type="button"
							onClick={openAuthModal}
							className="text-[#737373] text-sm hover:text-white transition-colors"
						>
							Complete your profile setup to join the discussion.
						</button>
					) : (
						<button
							type="button"
							onClick={openAuthModal}
							className="text-[#737373] text-sm hover:text-white transition-colors"
						>
							Sign in to join the discussion.
						</button>
					)}
				</div>

				{/* Post list */}
				{posts.length === 0 ? (
					<div className="p-8 text-center">
						<p className="text-[#404040] text-sm">
							No comments yet. Be the first.
						</p>
					</div>
				) : (
					<div className="divide-y divide-[#1a1a1a]">
						{posts.map((post) => (
							<div key={post.id} className="px-4 py-4">
								<div className="flex items-center gap-3 mb-2">
									<span className="font-grotesk font-semibold text-white text-xs">
										{post.user.username}
									</span>
									<span className="text-[#404040] text-xs font-mono">
										{new Date(post.created_at).toLocaleDateString("en-US", {
											year: "numeric",
											month: "short",
											day: "numeric",
										})}
									</span>
									{post.edited && (
										<span className="text-[#404040] text-[10px] uppercase tracking-widest">
											edited
										</span>
									)}
								</div>

								{editingId === post.id ? (
									<div>
										<textarea
											value={editBody}
											onChange={(e) => setEditBody(e.target.value)}
											maxLength={2000}
											rows={3}
											className="w-full bg-[#0a0a0a] border border-[#1a1a1a] text-white text-sm px-3 py-2 resize-none focus:outline-none focus:border-[#404040] transition-colors"
										/>
										<div className="flex gap-2 mt-2">
											<button
												type="button"
												onClick={() => handleEdit(post.id)}
												className="px-3 py-1 text-xs font-grotesk font-semibold uppercase tracking-widest bg-white text-black hover:bg-[#e0e0e0] transition-colors"
											>
												Save
											</button>
											<button
												type="button"
												onClick={() => setEditingId(null)}
												className="px-3 py-1 text-xs font-grotesk font-semibold uppercase tracking-widest text-[#737373] hover:text-white border border-[#1a1a1a] hover:border-[#404040] transition-colors"
											>
												Cancel
											</button>
										</div>
									</div>
								) : (
									<p className="text-[#a3a3a3] text-sm leading-relaxed whitespace-pre-wrap">
										{post.body}
									</p>
								)}

								{isOwner(post) && editingId !== post.id && (
									<div className="flex gap-3 mt-2">
										<button
											type="button"
											onClick={() => {
												setEditingId(post.id);
												setEditBody(post.body);
											}}
											className="text-[10px] uppercase tracking-widest text-[#404040] hover:text-[#737373] transition-colors"
										>
											Edit
										</button>
										<button
											type="button"
											onClick={() => handleDelete(post.id)}
											className="text-[10px] uppercase tracking-widest text-[#404040] hover:text-red-400 transition-colors"
										>
											Delete
										</button>
									</div>
								)}
							</div>
						))}
					</div>
				)}

				{/* Pagination */}
				{totalPages > 1 && (
					<div className="flex items-center justify-center gap-4 p-4 border-t border-[#1a1a1a]">
						<button
							type="button"
							disabled={page === 1}
							onClick={() => setPage((p) => p - 1)}
							className="text-xs uppercase tracking-widest text-[#737373] hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
						>
							Prev
						</button>
						<span className="text-xs font-mono text-[#404040]">
							{page} / {totalPages}
						</span>
						<button
							type="button"
							disabled={page === totalPages}
							onClick={() => setPage((p) => p + 1)}
							className="text-xs uppercase tracking-widest text-[#737373] hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
						>
							Next
						</button>
					</div>
				)}
			</div>
		</div>
	);
}
