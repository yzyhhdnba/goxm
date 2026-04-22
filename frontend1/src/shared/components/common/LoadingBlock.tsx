type LoadingBlockProps = {
  lines?: number;
};

export function LoadingBlock({ lines = 3 }: LoadingBlockProps) {
  return (
    <div className="animate-pulse rounded-md bg-white p-5 shadow-sm">
      {Array.from({ length: lines }).map((_, index) => (
        <div
          key={index}
          className={`mb-3 h-4 rounded-full bg-slate-200 ${index === lines - 1 ? 'mb-0 w-2/3' : 'w-full'}`}
        />
      ))}
    </div>
  );
}
