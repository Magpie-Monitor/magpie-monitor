import { Chart as ChartJS, ChartOptions, RadialLinearScale } from 'chart.js';
import { PolarArea } from 'react-chartjs-2';
import './PolarChart.scss';

ChartJS.register(RadialLinearScale);

export interface PolarChartEntry {
  value: number;
  label: string;
  color: string;
}

export interface PolarChartProps {
  data: PolarChartEntry[];
  label: string;
}

const PolarChart = ({ data, label }: PolarChartProps) => {
  const config = {
    labels: data.map((entry) => entry.label),
    datasets: [
      {
        data: data.map((entry) => entry.value),
        backgroundColor: data.map((entry) => entry.color),
        borderColor: 'transparent',
        borderWidth: 1,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: true,
    resizeDelay: 200,
    color: 'white',
    backgroundColor: 'none',
    plugins: {
      legend: {
        display: true,
        position: 'bottom',
      },
      title: {
        display: true,
        text: label,
        position: 'top',
        color: 'white',
        font: {
          size: 16,
        },
      },
    },
    scales: {
      r: {
        angleLines: {
          color: 'white',
        },
        grid: {},
        pointLabels: {
          font: {},
        },
        ticks: {
          backdropColor: 'transparent',
          font: {},
        },
      },
    },
  } as ChartOptions<'polarArea'>;

  return (
    <div className="polar-chart">
      <PolarArea data={config} options={options} />
    </div>
  );
};

export default PolarChart;
