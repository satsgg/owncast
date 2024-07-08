import { useContext, useEffect, useState } from 'react';
import { VideoQuality } from '../../../services/video-settings-service';
import { Button, Slider } from 'antd';
import styles from './L402.module.scss';
import { L402ServiceContext } from '../../../services/l402-service';
import { Lsat } from 'lsat-js';
import QRCode from 'react-qr-code';
import Paragraph from 'antd/lib/typography/Paragraph';

const DURATION_SLIDER_MARKS = {
  [2]: {
    style: {
      marginLeft: '24px',
    },
    label: `2 min`,
  },
  [60]: {
    style: {
      marginLeft: '-10px',
    },
    label: `1 hour`,
  },
};

const sleep = (milliseconds: number) => {
  return new Promise(resolve => setTimeout(resolve, milliseconds));
};

export const L402 = ({
  setL402,
  videoQualities,
}: {
  setL402: Function;
  videoQualities: VideoQuality[];
}) => {
  const L402Service = useContext(L402ServiceContext);
  // const [selectedQualityIndex, setSelectedQualityIndex] = useState<number | null>(null);
  const [selectedQuality, setSelectedQuality] = useState<VideoQuality | null>(null);
  const [modal, setModal] = useState<'quality' | 'duration' | 'payment' | 'none'>('quality');
  const [minutes, setMinutes] = useState<number>(1);
  const [price, setPrice] = useState<number>(0);
  const [l402Challenge, setL402Challenge] = useState<Lsat | null>(null);

  useEffect(() => {
    if (!l402Challenge) return;
    console.debug('l402.paymentPreimage', l402Challenge.paymentPreimage);
    const awaitPayment = async () => {
      let payment = null;
      while (true) {
        payment = await L402Service.checkPaymentStatus(l402Challenge.paymentHash);
        if (payment && payment.status === 'SETTLED') break;
        await sleep(1000);
      }
      const l402 = l402Challenge; // deep copy?
      l402.setPreimage(payment.preimage);
      setL402(l402);
      // setL402(l402Challenge.toToken());
    };
    awaitPayment();
  }, [l402Challenge]);

  const requestInvoice = async () => {
    // if (!selectedQualityIndex) return;
    if (!selectedQuality) return;
    // const l402Challenge = await L402Service.fetchInvoice(selectedQualityIndex, minutes * 60);
    const l402Challenge = await L402Service.fetchInvoice(selectedQuality.index, minutes * 60);
    if (!l402Challenge) {
      console.error('handle l402Challenge error in ui');
      return;
    }
    setL402Challenge(l402Challenge);
    setModal('payment');
  };

  return (
    <div className={styles.l402}>
      {
        {
          quality: (
            <div className={styles.l402Container}>
              <h1 style={{ textAlign: 'center' }}>Payment Required</h1>
              {videoQualities &&
                videoQualities.map((q, p) => (
                  <Button
                    key={p}
                    onClick={() => {
                      // setSelectedQualityIndex(p);
                      setSelectedQuality(q);
                      setPrice(60 * q.price);
                      setModal('duration');
                    }}
                  >
                    {q.name}, {q.price === 0 ? 'free' : q.price + ' millisats/sec'}
                  </Button>
                ))}
            </div>
          ),
          duration: (
            <div className={styles.l402Container}>
              <h1>Select Duration</h1>
              <p>
                Cost: {selectedQuality ? Math.ceil(selectedQuality.price / 1000) : 0} sats per
                second
              </p>
              <div className="segment-slider-container">
                <Slider
                  tipFormatter={value => `${value} minutes`}
                  disabled={false}
                  defaultValue={1}
                  value={minutes}
                  onChange={value => {
                    if (!selectedQuality) return;
                    setMinutes(value);
                    console.debug('value in minutes', value);
                    // setPrice(value * 60 * videoQualities[selectedQualityIndex].price);
                    setPrice(value * 60 * selectedQuality.price);
                  }}
                  step={1}
                  min={1}
                  max={60}
                  marks={DURATION_SLIDER_MARKS}
                />
              </div>
              {/* <p>Price: {calculatePrice(10, selectedQuality.price || 0)} sats</p> */}
              <p>Duration: {minutes} minutes</p>
              <p>Price: {Math.ceil(price / 1000)} sats</p>
              <Button onClick={() => setModal('quality')}>Go Back</Button>
              <Button type="primary" onClick={requestInvoice}>
                Purchase
              </Button>
            </div>
          ),
          payment: (
            <div className={styles.l402Container}>
              <h1 style={{ textAlign: 'center' }}>Invoice</h1>
              <div style={{ background: 'white', padding: '16px', borderRadius: '12px' }}>
                <QRCode value="hey" />
              </div>
              <Paragraph
                ellipsis
                style={{ maxWidth: '288px', color: 'white' }}
                copyable={{
                  text: l402Challenge ? l402Challenge.invoice : '',
                  onCopy: () => console.debug('copied'),
                }}
              >
                {l402Challenge ? l402Challenge.invoice : ''}
              </Paragraph>
              <Button
                onClick={() => {
                  setModal('none');
                }}
              >
                Pay with webln
              </Button>
              <Button
                onClick={() => {
                  setL402(null);
                  setModal('duration');
                }}
              >
                Go Back
              </Button>
            </div>
          ),
          none: null,
        }[modal]
      }
    </div>
  );
};
