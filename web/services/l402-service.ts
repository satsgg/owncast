import { Lsat } from 'lsat-js';
import { createContext } from 'react';

export interface L402StaticService {
  fetchInvoice(quality_index: number, seconds: number): Promise<Lsat>;
  checkPaymentStatus(hash: string): Promise<any>;
}

class L402Service {
  // TODO: Add price to response
  private static readonly STREAM_BASE_URL = '/hls';
  private static readonly STREAM_PLAYLIST_FILE = 'stream.m3u8';

  public static async fetchInvoice(quality_index: number, seconds: number): Promise<Lsat> {
    try {
      const response = await fetch(
        `${L402Service.STREAM_BASE_URL}/${quality_index}/${this.STREAM_PLAYLIST_FILE}?t=${seconds}`,
      );
      if (response.status !== 402) throw new Error('Some other shit happened');

      console.debug('WWW-Authenticate', response.headers.get('WWW-Authenticate'));
      const l402 = Lsat.fromHeader(response.headers.get('WWW-Authenticate'));

      return l402;
    } catch (e) {
      console.error(e);
    }
    return null;
  }

  private static readonly PAYMENT_STATUS_PATH = '/.well-known/bolt11';

  // r.HandleFunc("/.well-known/bolt11", controllers.CheckPaymentStatus)
  public static async checkPaymentStatus(hash: string): Promise<any> {
    try {
      const response = await fetch(`${this.PAYMENT_STATUS_PATH}?h=${hash}`);
      console.debug('response', response);
      const status = await response.json();
      console.debug('status (json)', status);
      return status;
    } catch (e) {
      console.error(e);
    }
  }
}

export const L402ServiceContext = createContext<L402StaticService>(L402Service);
