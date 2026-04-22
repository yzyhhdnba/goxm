import { ChangeEvent, FormEvent, useEffect, useMemo, useRef, useState } from 'react';
import { FeedAPI } from '@/api/modules/feed';
import { UploadAPI } from '@/api/modules/upload';
import { AppShell } from '@/shared/components/layout/AppShell';
import { EmptyState } from '@/shared/components/common/EmptyState';
import { LoadingBlock } from '@/shared/components/common/LoadingBlock';
import { formatDate, resolveMediaUrl } from '@/shared/utils/format';
import type { Area, CreatorVideoItem } from '@/shared/types/domain';

type CreatorTab = 'pending' | 'approved' | 'rejected' | 'all';

const pageSize = 20;

export function CreatorStudioPage() {
  const [areas, setAreas] = useState<Area[]>([]);
  const [videos, setVideos] = useState<CreatorVideoItem[]>([]);
  const [activeTab, setActiveTab] = useState<CreatorTab>('pending');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [areaId, setAreaId] = useState<number>(0);
  const [videoFile, setVideoFile] = useState<File | null>(null);
  const [coverFile, setCoverFile] = useState<File | null>(null);
  const [videoPreviewUrl, setVideoPreviewUrl] = useState('');
  const [coverPreviewUrl, setCoverPreviewUrl] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [loadingList, setLoadingList] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const videoRef = useRef<HTMLVideoElement | null>(null);

  useEffect(() => {
    async function loadAreas() {
      try {
        const response = await FeedAPI.getAreas();
        setAreas(response);
        if (response[0]) {
          setAreaId((current) => current || response[0].id);
        }
      } catch (requestError) {
        setError('分区加载失败');
      }
    }

    void loadAreas();
  }, []);

  useEffect(() => {
    async function loadList() {
      setLoadingList(true);
      try {
        const response = await UploadAPI.listCreatorVideos({
          review_status: activeTab,
          page: 1,
          page_size: pageSize,
        });
        setVideos(response.list);
      } catch (requestError) {
        setError('稿件列表加载失败');
      } finally {
        setLoadingList(false);
      }
    }

    void loadList();
  }, [activeTab]);

  const canSubmit = useMemo(() => title.trim() && description.trim() && areaId && videoFile && coverFile, [areaId, coverFile, description, title, videoFile]);

  const onVideoChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] || null;
    setVideoFile(file);
    if (file) {
      setVideoPreviewUrl(URL.createObjectURL(file));
    } else {
      setVideoPreviewUrl('');
    }
  };

  const onCoverChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] || null;
    setCoverFile(file);
    if (file) {
      setCoverPreviewUrl(URL.createObjectURL(file));
    } else {
      setCoverPreviewUrl('');
    }
  };

  const captureCover = async () => {
    if (!videoRef.current) {
      return;
    }
    const video = videoRef.current;
    const canvas = document.createElement('canvas');
    canvas.width = video.videoWidth || 640;
    canvas.height = video.videoHeight || 360;
    const context = canvas.getContext('2d');
    if (!context) {
      return;
    }
    context.drawImage(video, 0, 0, canvas.width, canvas.height);
    const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/jpeg', 0.92));
    if (!blob) {
      return;
    }
    const file = new File([blob], 'captured-cover.jpg', { type: 'image/jpeg' });
    setCoverFile(file);
    setCoverPreviewUrl(URL.createObjectURL(file));
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!canSubmit || !videoFile || !coverFile) {
      return;
    }

    setSubmitting(true);
    setError(null);
    try {
      const created = await UploadAPI.createVideo({
        area_id: areaId,
        title: title.trim(),
        description: description.trim(),
      });
      await UploadAPI.uploadSource(created.id, videoFile);
      await UploadAPI.uploadCover(created.id, coverFile);
      setTitle('');
      setDescription('');
      setVideoFile(null);
      setCoverFile(null);
      setVideoPreviewUrl('');
      setCoverPreviewUrl('');
      setActiveTab('pending');
      const response = await UploadAPI.listCreatorVideos({ review_status: 'pending', page: 1, page_size: pageSize });
      setVideos(response.list);
    } catch (requestError) {
      setError('投稿失败，请检查文件和表单内容');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <AppShell>
      <section className="mb-6 rounded-md bg-white p-6 shadow-sm">
        <p className="text-sm uppercase tracking-[0.3em] text-sea">Creator Studio</p>
        <h1 className="mt-2 text-3xl font-semibold text-ink">视频上传发布</h1>
        <p className="mt-3 text-sm text-slate-500">先提交稿件元数据，再上传原生视频文件，封面可以手动上传，也可以从预览视频当前帧截取。</p>
      </section>

      {error ? <div className="mb-6 rounded-md bg-amber-50 px-4 py-3 text-sm text-amber-700">{error}</div> : null}

      <div className="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
        <form className="space-y-5 rounded-md bg-white p-6 shadow-sm" onSubmit={handleSubmit}>
          <div>
            <label className="mb-2 block text-sm font-medium text-slate-600">标题</label>
            <input className="w-full rounded-md border border-slate-200 bg-slate-50 px-4 py-3 outline-none focus:border-accent" onChange={(e) => setTitle(e.target.value)} value={title} />
          </div>
          <div>
            <label className="mb-2 block text-sm font-medium text-slate-600">分区</label>
            <select className="w-full rounded-md border border-slate-200 bg-slate-50 px-4 py-3 outline-none focus:border-accent" onChange={(e) => setAreaId(Number(e.target.value))} value={areaId}>
              {areas.map((area) => (
                <option key={area.id} value={area.id}>
                  {area.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="mb-2 block text-sm font-medium text-slate-600">简介</label>
            <textarea className="min-h-32 w-full rounded-md border border-slate-200 bg-slate-50 px-4 py-3 outline-none focus:border-accent" onChange={(e) => setDescription(e.target.value)} value={description} />
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            <label className="rounded-md border border-dashed border-slate-300 bg-paper px-4 py-6 text-sm text-slate-600">
              选择原生视频
              <input accept="video/*" className="mt-3 block w-full text-xs" onChange={onVideoChange} type="file" />
            </label>
            <label className="rounded-md border border-dashed border-slate-300 bg-paper px-4 py-6 text-sm text-slate-600">
              选择封面图片
              <input accept="image/*" className="mt-3 block w-full text-xs" onChange={onCoverChange} type="file" />
            </label>
          </div>
          <button className="rounded-full bg-ink px-5 py-3 text-sm font-medium text-white disabled:cursor-not-allowed disabled:opacity-60" disabled={!canSubmit || submitting} type="submit">
            {submitting ? '上传中...' : '提交稿件'}
          </button>
        </form>

        <section className="space-y-5 rounded-md bg-white p-6 shadow-sm">
          <div>
            <p className="text-sm uppercase tracking-[0.25em] text-slate-400">Preview</p>
            <h2 className="text-2xl font-semibold text-ink">视频与封面预览</h2>
          </div>
          {videoPreviewUrl ? (
            <div className="space-y-3">
              <video className="aspect-video w-full rounded-md bg-ink" controls ref={videoRef} src={videoPreviewUrl} />
              <button className="rounded-full border border-slate-200 px-4 py-2 text-sm text-slate-600" onClick={captureCover} type="button">
                截取当前帧作为封面
              </button>
            </div>
          ) : (
            <EmptyState title="先选择视频文件" />
          )}
          {coverPreviewUrl ? <img alt="cover preview" className="aspect-video w-full rounded-md object-cover" src={coverPreviewUrl} /> : null}
        </section>
      </div>

      <section className="mt-6 space-y-5">
        <div className="flex flex-wrap gap-3">
          {(['pending', 'approved', 'rejected', 'all'] as CreatorTab[]).map((tab) => (
            <button
              className={`rounded-full px-4 py-2 text-sm transition ${activeTab === tab ? 'bg-ink text-white' : 'bg-white text-slate-600'}`}
              key={tab}
              onClick={() => setActiveTab(tab)}
              type="button"
            >
              {tab === 'pending' ? '待审核' : tab === 'approved' ? '已通过' : tab === 'rejected' ? '未通过' : '全部'}
            </button>
          ))}
        </div>
        {loadingList ? (
          <div className="space-y-4">
            {Array.from({ length: 3 }).map((_, index) => (
              <LoadingBlock key={index} lines={5} />
            ))}
          </div>
        ) : videos.length === 0 ? (
          <EmptyState title="当前没有对应状态的稿件" />
        ) : (
          <div className="space-y-4">
            {videos.map((item) => (
              <article className="rounded-md bg-white p-5 shadow-sm" key={item.id}>
                <div className="flex flex-col gap-4 md:flex-row">
                  <img alt={item.title} className="aspect-video w-full rounded-md object-cover md:w-72" src={resolveMediaUrl(item.cover_url)} />
                  <div className="min-w-0 flex-1 space-y-3">
                    <div className="flex flex-wrap items-center gap-3">
                      <h3 className="text-xl font-semibold text-ink">{item.title}</h3>
                      <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">{item.area_name}</span>
                      <span className="rounded-full bg-paper px-3 py-1 text-xs text-slate-500">{item.review_status}</span>
                    </div>
                    <p className="text-sm text-slate-600">{item.description}</p>
                    {item.review_reason ? <p className="text-sm text-rose-500">驳回原因：{item.review_reason}</p> : null}
                    <p className="text-xs text-slate-400">更新时间：{formatDate(item.updated_at)}</p>
                  </div>
                </div>
              </article>
            ))}
          </div>
        )}
      </section>
    </AppShell>
  );
}
