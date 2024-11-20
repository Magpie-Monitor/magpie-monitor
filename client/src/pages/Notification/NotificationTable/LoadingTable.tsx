
export interface LoadingTableProps {
  isLoading: boolean;
  children: React.ReactNode;
}

const LoadingTable = ({ isLoading, children }: LoadingTableProps) => {
  return isLoading ? <p>Loading...</p> : children;
};

export default LoadingTable;