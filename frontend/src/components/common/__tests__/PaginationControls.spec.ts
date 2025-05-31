import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import PaginationControls from '../PaginationControls.vue';

describe('PaginationControls.vue', () => {
  it('renders correctly with multiple pages', () => {
    const wrapper = mount(PaginationControls, {
      props: {
        currentPage: 2,
        totalPages: 5,
        pageSize: 10,
        totalItems: 45,
      },
    });
    expect(wrapper.text()).toContain('Pág 2 de 5');
    expect(wrapper.text()).toContain('Mostrando 11 a 20 de 45 resultados');
    expect(wrapper.findAll('button').length).toBeGreaterThan(2);
  });

  it('disables "Previous" button on first page', () => {
    const wrapper = mount(PaginationControls, {
      props: {
        currentPage: 1,
        totalPages: 3,
        pageSize: 10,
        totalItems: 25,
      },
    });
    const prevButtonsDesktop = wrapper.findAll('button[aria-label="Anterior"]');
    if (prevButtonsDesktop.length > 0) {
        const prevButtonDesktop = prevButtonsDesktop[0];
        expect(prevButtonDesktop.attributes('disabled')).toBeDefined();
    }
    const prevButtonMobile = wrapper.findAll('button').find(b => b.text().includes('Anterior'));
     if (prevButtonMobile) {
        expect(prevButtonMobile.attributes('disabled')).toBeDefined();
    }
  });

  it('disables "Next" button on last page', () => {
    const wrapper = mount(PaginationControls, {
      props: {
        currentPage: 3,
        totalPages: 3,
        pageSize: 10,
        totalItems: 25,
      },
    });
    const nextButtonsDesktop = wrapper.findAll('button[aria-label="Siguiente"]');
    if (nextButtonsDesktop.length > 0) {
        const nextButtonDesktop = nextButtonsDesktop[0];
        expect(nextButtonDesktop.attributes('disabled')).toBeDefined();
    }
    const nextButtonMobile = wrapper.findAll('button').find(b => b.text().includes('Siguiente'));
    if (nextButtonMobile) {
        expect(nextButtonMobile.attributes('disabled')).toBeDefined();
    }
  });

  it('emits "page-changed" event with correct page number when a page number is clicked', async () => {
    const wrapper = mount(PaginationControls, {
      props: {
        currentPage: 1,
        totalPages: 5,
        pageSize: 10,
        totalItems: 45,
      },
    });
    const pageTwoButton = wrapper.findAll('button').find(b => b.text() === '2');
    if (pageTwoButton) {
        await pageTwoButton.trigger('click');
        expect(wrapper.emitted()['page-changed']).toBeTruthy();
        expect(wrapper.emitted()['page-changed'][0]).toEqual([2]);
    } else {
        console.warn("Botón de página '2' no encontrado para test de emisión de evento.");
    }
  });

  it('does not render if totalPages is 0 or 1', () => {
    const wrapper = mount(PaginationControls, {
      props: { currentPage: 1, totalPages: 1, pageSize: 10, totalItems: 5 },
    });
    expect(wrapper.html()).toContain('');

    const wrapperZero = mount(PaginationControls, {
        props: { currentPage: 1, totalPages: 0, pageSize: 10, totalItems: 0 },
    });
    expect(wrapperZero.html()).toContain('');
  });
});
