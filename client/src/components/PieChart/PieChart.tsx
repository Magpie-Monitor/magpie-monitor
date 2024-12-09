import { Chart as ChartJS, ChartOptions, RadialLinearScale } from 'chart.js';
import { Pie } from 'react-chartjs-2';
import './PieChart.scss';

ChartJS.register(RadialLinearScale);

export interface PieChartEntry {
  value: number;
  label: string;
  color: string;
}

export interface PieChartProps {
  data: PieChartEntry[];
  label: string;
}

const PieChart = ({ data, label }: PieChartProps) => {
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
    maintainAspectRatio: false,
    resizeDelay: 200,
    color: 'white',
    backgroundColor: 'none',
    plugins: {
      legend: {
        display: true,
        position: 'bottom',
        labels: {
          font: {
            family: 'Roboto',
            weight: 'bold',
          },
          color: 'white',
        },
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
  } as ChartOptions<'pie'>;

  return (
    <div className="pie-chart">
      <Pie data={config} options={options} />
    </div>
  );
};

export default PieChart;
