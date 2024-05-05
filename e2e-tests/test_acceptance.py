from os.path import join
import re
from playwright.sync_api import Page, Browser, expect
from fatetable.pagemodel import IndexScene, CreateSessionModal, GameMasterScene, PlayerScene, AddAspectModal, JoinSessionModal

timeout=15000

def test_index_page_with_roll_buttons(page: Page):
    index_scene = IndexScene(page)
    index_scene.visit_page()

    # Expect a title to be correct
    expect(page).to_have_title("Fate Core Remote Table")

    expect(index_scene.skill_check_result()).not_to_be_visible()
    # Expect dice roll buttons to be visible
    for i in range(5):
        expect(index_scene.skill_check_btn(i)).to_be_visible()
        index_scene.skill_check_btn(i).click()
        expect(index_scene.skill_check_result()).to_be_visible()

    expect(index_scene.join_session_btn()).to_be_visible()
    expect(index_scene.create_new_session_btn()).to_be_visible()


def test_full_game_scenario(browser: Browser, base_url: str):
    gm_context = browser.new_context(base_url=base_url)
    gm_page = gm_context.new_page()
    gm_page.on("console", lambda msg: print(f"GM page: {msg.text}"))

    player_one_context = browser.new_context(base_url=base_url)
    player_one_page = player_one_context.new_page()
    player_one_page.on(
        "console", lambda msg: print(f"P1 page: {msg.text}"))

    try:
        #########################################
        # GM -> Visit HomePage, start new session
        gm_index_scene = IndexScene(gm_page)
        gm_index_scene.visit_page()

        expect(gm_index_scene.create_new_session_btn()).to_be_visible()
        gm_index_scene.create_new_session_btn().click()

        create_new_session_modal = CreateSessionModal(gm_page)
        expect(create_new_session_modal.locator).to_be_visible()

        create_new_session_modal.session_title_input().fill("Test Session")
        create_new_session_modal.ok_btn().click()
        expect(create_new_session_modal.locator).to_be_hidden()

        #########################################
        # GM -> See an empty game master scene

        gm_scene = GameMasterScene(gm_page)
        expect(gm_scene.title()).to_have_text(re.compile(r'^[A-Z]{2} @ Test Session'))
        expect(gm_scene.add_aspect_btn()).to_be_visible()
        expect(gm_scene.join_session_link()).to_be_visible()

        session_id = gm_page.url.split('/')[-1]

        #########################################
        # Player One -> Visit HomePage, join started session
        player_one_index_scene = IndexScene(player_one_page)
        player_one_index_scene.visit_page()
        player_one_index_scene.join_session_btn().click()

        join_session_modal = JoinSessionModal(player_one_page)
        expect(join_session_modal.locator).to_be_visible()
        join_session_modal.session_id_input().fill(session_id)
        join_session_modal.player_name_input().fill("Player One")
        join_session_modal.ok_btn().click()

        player_one_scene = PlayerScene(player_one_page)
        expect(player_one_scene.fate_points()).to_have_text('0')
        expect(player_one_scene.spend_fate_points_btn()).to_be_disabled()

        #########################################
        # GM:
        # - Wait for player to be visible, then
        # - increase fate points for player one
        # - add global aspect
        expect(gm_scene.player_locator(0)).to_be_visible(timeout=timeout)
        player_one_card = gm_scene.player(0)
        expect(player_one_card.fate_points_label()).to_have_text('0')
        expect(player_one_card.dec_fate_points_btn()).to_be_disabled()
        expect(player_one_card.inc_fate_points_btn()).to_be_enabled()

        player_one_card.inc_fate_points_btn().click()
        expect(player_one_card.fate_points_label()).to_have_text('1', timeout=timeout)

        gm_scene.add_aspect_btn().click()
        add_aspect_modal = AddAspectModal(gm_page)
        add_aspect_modal.aspect_name_input().fill("Foggy")
        add_aspect_modal.ok_btn().click()

        #########################################
        # Player One:
        # - Wait for aspect and fate points to be visible, then
        # - spend fate point
        expect(player_one_scene.fate_points()).to_have_text('1', timeout=timeout)
        player_one_scene.spend_fate_points_btn().click()
        expect(player_one_scene.fate_points()).to_have_text('0', timeout=timeout)

        #########################################
        # GM:
        # - Wait for player one fate points to drop to zero
        expect(player_one_card.fate_points_label()).to_have_text('0', timeout=timeout)

    finally:
        player_one_page.screenshot(path=join('test-results', 'player.png'))
        gm_page.screenshot(path=join('test-results', 'gm.png'))
        
        player_one_context.close()
        gm_context.close()

